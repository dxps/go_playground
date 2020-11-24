package server

import (
	"context"

	api "devisions.org/go-dist-svcs/log/api/v1"
	"google.golang.org/grpc"
)

// This line is a trick which guarantees that *grpcServer type satisfies the `api.LogServer` interface.
var _ api.LogServer = (*grpcServer)(nil)

type Config struct {
	CommitLog CommitLog
}

func NewGRPCServer(config *Config) (*grpc.Server, error) {

	gsrv := grpc.NewServer()
	srv, err := newgrpcServer(config)
	if err != nil {
		return nil, err
	}
	api.RegisterLogServer(gsrv, srv)
	return gsrv, nil
}

type grpcServer struct {
	*Config
}

func newgrpcServer(config *Config) (srv *grpcServer, err error) {
	srv = &grpcServer{
		Config: config,
	}
	return srv, nil
}

func (s *grpcServer) Produce(ctx context.Context, req *api.ProduceRequest) (*api.ProduceResponse, error) {

	offset, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &api.ProduceResponse{Offset: offset}, nil
}

func (s *grpcServer) Consume(ctx context.Context, req *api.ConsumeRequest) (*api.ConsumeResponse, error) {

	record, err := s.CommitLog.Read(req.Offset)
	if err != nil {
		return nil, err
	}
	return &api.ConsumeResponse{Record: record}, nil
}

// ProduceStream implements a bidirectional streaming RPC: the client can stream data into server's log
// and the server can tell the client whether each request succeeded.
func (s *grpcServer) ProduceStream(stream api.Log_ProduceStreamServer) error {

	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		res, err := s.Produce(stream.Context(), req)
		if err != nil {
			return err
		}
		if err = stream.Send(res); err != nil {
			return err
		}
	}
}

// ConsumeStream implements a server-side streaming RPC: the client can tell the server
// where in the log to read records, and the server will stream every record that follows
// (including any future records that may appear).
func (s *grpcServer) ConsumeStream(req *api.ConsumeRequest, stream api.Log_ConsumeStreamServer) error {

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			// The server reads to the end of the log and after that it loops
			// until another record is produced that it also sends to the client.
			res, err := s.Consume(stream.Context(), req)
			switch err.(type) {
			case nil:
			case api.ErrOffsetOutOfRange:
				continue
			default:
				return err
			}
			if err = stream.Send(res); err != nil {
				return err
			}
			req.Offset++
		}
	}
}

type CommitLog interface {
	Append(*api.Record) (uint64, error)
	Read(uint64) (*api.Record, error)
}
