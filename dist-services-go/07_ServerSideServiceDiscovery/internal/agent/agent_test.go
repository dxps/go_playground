package agent_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	api "devisions.org/go-dist-svcs/log/api/v1"
	"devisions.org/go-dist-svcs/log/internal/agent"
	"devisions.org/go-dist-svcs/log/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestAgent(t *testing.T) {

	serverTLSConfig, err := config.SetupTLSConfig(
		config.TLSConfig{
			CertFile:      config.ServerCertFile,
			KeyFile:       config.ServerKeyFile,
			CAFile:        config.CAFile,
			ServerAddress: "127.0.0.1",
			Server:        true,
		})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetupTLSConfig(
		config.TLSConfig{
			CertFile:      config.RootClientCertFile,
			KeyFile:       config.RootClientKeyFile,
			CAFile:        config.CAFile,
			ServerAddress: "127.0.0.1",
			Server:        false,
		})
	require.NoError(t, err)

	// Setting up a three node cluster.
	var agents []*agent.Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]

		dataDir, err := ioutil.TempDir("", "server-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		agent, err := agent.New(agent.Config{
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			DataDir:         dataDir,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			NodeName:        fmt.Sprintf("node%d", i),
			StartJoinAddrs:  startJoinAddrs,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
		})
		require.NoError(t, err)

		log.Printf("\n\n[TEST] Agent created: %+v\n\n", agent)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err = agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()

	log.Printf("\n\n[TEST] %d agents were started.\n\n", len(agents))
	// giving time for the cluster to get ready (agents to start and discover each other)
	time.Sleep(3 * time.Second)

	// Testing the produce and consume from a single node.
	leaderClient := client(t, agents[0], peerTLSConfig)

	produceResp, err := leaderClient.Produce(
		context.Background(),
		&api.ProduceRequest{Record: &api.Record{Value: []byte("foo")}},
	)
	require.NoError(t, err)

	consumeResp, err := leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{Offset: produceResp.Offset},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResp.Record.Value, []byte("foo"))

	log.Printf("\n\n[TEST] Leader client produced then consumed %+v\n\n", consumeResp)

	// Testing that another node replicated the record.
	// First, let's give some time for the replication to happen.
	time.Sleep(3 * time.Second)

	log.Printf("\n\n[TEST] Follower client is trying to consume from agent 1 (%+v)\n\n", agents[1])

	followerClient := client(t, agents[1], peerTLSConfig)
	consumeResp, err = followerClient.Consume(
		context.Background(),
		&api.ConsumeRequest{Offset: produceResp.Offset},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResp.Record.Value, []byte("foo"))
}

func client(t *testing.T, agent *agent.Agent, tlsConfig *tls.Config) api.LogClient {

	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	cc, err := grpc.Dial(rpcAddr, opts...)
	require.NoError(t, err)

	return api.NewLogClient(cc)
}
