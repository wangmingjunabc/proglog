package loadbalance

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/wangmingjunabc/proglog/api/v1"
	"github.com/wangmingjunabc/proglog/internal/config"
	"github.com/wangmingjunabc/proglog/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestResolver(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	serverCreds := credentials.NewTLS(tlsConfig)
	srv, err := server.NewGRPCServer(&server.Config{GetServer: &getServers{}}, grpc.Creds(serverCreds))
	require.NoError(t, err)
	go srv.Serve(ln)

	conn := &clientConn{}
	tlsConfig, err = config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootCertFile,
		KeyFile:       config.RootKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: "127.0.0.1",
		Server:        false,
	})
	require.NoError(t, err)
	clientCreds := credentials.NewTLS(tlsConfig)
	opts := resolver.BuildOptions{
		DialCreds: clientCreds,
	}
	r := &Resolver{}
	_, err = r.Build(resolver.Target{
		Endpoint: ln.Addr().String(),
	}, conn, opts)
	require.NoError(t, err)

	wantState := resolver.State{
		Addresses: []resolver.Address{
			{Addr: "localhost:9001", Attributes: attributes.New("is_leader", true)},
			{Addr: "localhost:9002", Attributes: attributes.New("is_leader", false)},
		},
	}
	require.Equal(t, wantState, conn.state)
	conn.state.Addresses = nil
	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

type getServers struct{}

func (s *getServers) GetServers() ([]*api.Server, error) {
	return []*api.Server{
		{Id: "leader", RpcAddr: "localhost:9001", IsLeader: true},
		{Id: "follower", RpcAddr: "localhost:9002", IsLeader: false},
	}, nil
}

type clientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (c *clientConn) UpdateState(state resolver.State) error {
	c.state = state
	return nil
}

func (c *clientConn) ReportError(err error) {}

func (c *clientConn) NewAddress(addrs []resolver.Address) {}

func (c *clientConn) NewServiceConfig(config string) {}

func (c *clientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}
