package arpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/arpc"
)

func TestHealthServerOK(t *testing.T) {
	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return nil },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep1"},
			"service2": {},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthStatus, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, healthStatus.GetStatus())

	healthStatus, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service1"})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, healthStatus.GetStatus())

	healthStatus, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service2"})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, healthStatus.GetStatus())
}

func TestHealthServerUnknownService(t *testing.T) {
	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return nil },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep1"},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "foo"})
	testutils.RequireGRPCCodesEqual(t, err, codes.NotFound)
}

func TestHealthServerUnknownDependency(t *testing.T) {
	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return nil },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep3"},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service1"})
	testutils.RequireGRPCCodesEqual(t, err, codes.Unknown)
}

func TestHealthServerKO(t *testing.T) {
	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return nil },
			"dep2": func() error { return errors.New("uwups") },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep1"},
			"service2": {"dep1", "dep2"},
			"service3": {"dep2"},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthStatus, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, healthStatus.GetStatus())

	healthStatus, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service1"})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, healthStatus.GetStatus())

	healthStatus, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service2"})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, healthStatus.GetStatus())

	healthStatus, err = client.Check(ctx, &healthpb.HealthCheckRequest{Service: "service3"})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, healthStatus.GetStatus())
}

func TestHealthServerWatch(t *testing.T) {
	errs := map[string]error{
		"dep1": nil,
		"dep2": nil,
	}

	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return errs["dep1"] },
			"dep2": func() error { return errs["dep2"] },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep1"},
			"service2": {"dep1", "dep2"},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	ctxAll, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxService1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctxService2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)

	streamAll, err := client.Watch(ctxAll, &healthpb.HealthCheckRequest{Service: ""})
	require.NoError(t, err)
	streamService1, err := client.Watch(ctxService1, &healthpb.HealthCheckRequest{Service: "service1"})
	require.NoError(t, err)
	streamService2, err := client.Watch(ctxService2, &healthpb.HealthCheckRequest{Service: "service2"})
	require.NoError(t, err)

	// Every stream should report OK.
	statusAll, err := streamAll.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusAll.GetStatus())

	statusService1, err := streamService1.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusService1.GetStatus())

	statusService2, err := streamService2.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusService2.GetStatus())

	// Now we break dep2.
	errs["dep2"] = errors.New("uwups")

	statusAll, err = streamAll.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, statusAll.GetStatus())

	statusService1, err = streamService1.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusService1.GetStatus())

	statusService2, err = streamService2.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_NOT_SERVING, statusService2.GetStatus())

	// Back to normal.
	errs["dep2"] = nil

	statusAll, err = streamAll.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusAll.GetStatus())

	statusService1, err = streamService1.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusService1.GetStatus())

	statusService2, err = streamService2.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, statusService2.GetStatus())
}

func TestHealthServerWatchError(t *testing.T) {
	depsCheck := &arpc.DepsCheck{
		Dependencies: arpc.DepCheckCallbacks{
			"dep1": func() error { return nil },
		},
		Services: arpc.DepCheckServices{
			"service1": {"dep1"},
		},
	}

	listener, server, err := arpc.StartServer(8080)
	require.NoError(t, err)
	defer arpc.CloseServer(listener, server)

	healthServer := arpc.NewHealthServer(depsCheck, 100*time.Millisecond)
	healthpb.RegisterHealthServer(server, healthServer)

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	go func() {
		require.NoError(t, server.Serve(listener))
	}()

	ctxAll, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := healthpb.NewHealthClient(conn)

	streamAll, err := client.Watch(ctxAll, &healthpb.HealthCheckRequest{Service: "foo"})
	require.NoError(t, err)

	// Every stream should report OK.
	_, err = streamAll.Recv()
	testutils.RequireGRPCCodesEqual(t, err, codes.Canceled)
}
