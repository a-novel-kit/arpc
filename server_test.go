package arpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"

	"github.com/a-novel-kit/arpc"
	arpcmocks "github.com/a-novel-kit/arpc/mocks"
)

func setupServerStubServer(t *testing.T) *arpcmocks.StubServer {
	t.Helper()

	return &arpcmocks.StubServer{
		EmptyCallF: func(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
			return new(testgrpc.Empty), nil
		},
	}
}

func TestServerServing(t *testing.T) {
	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	testgrpc.RegisterTestServiceServer(server, setupServerStubServer(t))

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	require.NoError(t, err)
}

func TestServerNoPort(t *testing.T) {
	_, _, err := arpc.StartServer(0)
	require.ErrorIs(t, err, arpc.ErrPortRequired)
}
