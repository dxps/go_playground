package internal

import (
	"context"
	"log"
	"net"

	pb "github.com/devisions/go-playground/go-grpc-ur/s03-order-mgmt-svc/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	orders map[string]*pb.Order
}

func StartServer(address string) (*grpc.Server, error) {

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.Wrap(err, "trying to start the listener")
	}
	s := server{}
	gs := grpc.NewServer()
	pb.RegisterOrderMgmtServer(gs, &s)
	log.Println(">>> Starting the server and listening on", address)
	if err := gs.Serve(lis); err != nil {
		return nil, errors.Wrap(err, "trying to start the gRPC server")
	}
	return gs, nil
}

func (s *server) GetOrder(ctx context.Context, id *wrapperspb.StringValue) (*pb.Order, error) {

	ord, exists := s.orders[id.Value]
	if !exists {
		return nil, status.Error(codes.NotFound, "")
	}
	return ord, nil
}
