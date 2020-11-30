package agent

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	api "devisions.org/go-dist-svcs/log/api/v1"
	"devisions.org/go-dist-svcs/log/internal/auth"
	"devisions.org/go-dist-svcs/log/internal/discovery"
	"devisions.org/go-dist-svcs/log/internal/log"
	"devisions.org/go-dist-svcs/log/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Agent struct {
	Config

	log        *log.Log
	server     *grpc.Server
	membership *discovery.Membership
	replicator *log.Replicator

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

type Config struct {
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	DataDir         string
	BindAddr        string
	RPCPort         int
	NodeName        string
	StartJoinAddrs  []string
	ACLModelFile    string
	ACLPolicyFile   string
}

// RPCAddr returns the `IP:port` pair of IP address from `BindAddr` and port from `RPCPort`.
func (c Config) RPCAddr() (string, error) {

	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

func New(config Config) (*Agent, error) {

	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}
	setupFns := []func() error{
		a.setupLog,
		a.setupServer,
		a.setupMembership,
	}
	for _, fn := range setupFns {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func (a *Agent) setupLog() error {

	var err error
	a.log, err = log.NewLog(a.Config.DataDir, log.Config{})
	return err
}

func (a *Agent) setupServer() error {

	authzr := auth.New(a.Config.ACLModelFile, a.Config.ACLPolicyFile)
	srvCfg := &server.Config{CommitLog: a.log, Authorizer: authzr}

	var opts []grpc.ServerOption
	if a.Config.ServerTLSConfig != nil {
		creds := credentials.NewTLS(a.Config.ServerTLSConfig)
		opts = append(opts, grpc.Creds(creds))
	}
	var err error
	a.server, err = server.NewGRPCServer(srvCfg, opts...)
	if err != nil {
		return err
	}
	rpcAddr, err := a.RPCAddr()
	if err != nil {
		return err
	}
	lsnr, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}
	go func() {
		if err := a.server.Serve(lsnr); err != nil {
			fmt.Printf("Agent's server unable to listen to GRPC address %v due to: %s", rpcAddr, err)
			_ = a.Shutdown()
		}
	}()
	return err
}

// Internal setup of the `Replicator` and `Membership`.
func (a *Agent) setupMembership() error {

	rpcAddr, err := a.Config.RPCAddr()
	if err != nil {
		return err
	}
	var opts []grpc.DialOption
	if a.Config.PeerTLSConfig != nil {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(a.Config.PeerTLSConfig)))
	}
	conn, err := grpc.Dial(rpcAddr, opts...)
	if err != nil {
		return err
	}
	client := api.NewLogClient(conn)
	a.replicator = &log.Replicator{
		DialOpts:    opts,
		LocalServer: client,
	}
	a.membership, err = discovery.New(a.replicator, discovery.Config{
		NodeName: a.Config.NodeName,
		BindAddr: a.Config.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: a.Config.StartJoinAddrs,
	})
	return err
}

func (a *Agent) Shutdown() error {

	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()
	if a.shutdown {
		return nil
	}
	a.shutdown = true
	close(a.shutdowns)

	shutdownFns := []func() error{
		a.membership.Leave,
		a.replicator.Close,
		func() error {
			a.server.GracefulStop()
			return nil
		},
		a.log.Close,
	}
	for _, fn := range shutdownFns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
