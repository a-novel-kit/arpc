package arpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/arpc"
	arpcmocks "github.com/a-novel-kit/arpc/mocks"
	x509mocks "github.com/a-novel-kit/arpc/mocks/x509/x509"
)

type stubServerParams struct {
	insecure bool

	authority string
	token     *oauth2.Token
}

func setupClientStubServer(t *testing.T, params stubServerParams) *arpcmocks.StubServer {
	t.Helper()

	return &arpcmocks.StubServer{
		EmptyCallF: func(ctx context.Context, _ *testgrpc.Empty) (*testgrpc.Empty, error) {
			pr, ok := peer.FromContext(ctx)
			if !ok {
				return nil, status.Error(codes.DataLoss, "Failed to get peer from ctx")
			}

			expectedSecLevel := lo.Ternary(params.insecure, credentials.NoSecurity, credentials.PrivacyAndIntegrity)
			if err := credentials.CheckSecurityLevel(pr.AuthInfo, expectedSecLevel); err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "Wrong security level: %s", err)
			}

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.DataLoss, "Failed to get metadata from ctx")
			}

			contentType, ok := md["content-type"]
			if !ok || len(contentType) == 0 {
				return nil, status.Error(codes.DataLoss, "Failed to get content type from metadata")
			}
			if contentType[0] != "application/grpc" {
				return nil, status.Errorf(codes.InvalidArgument, "Wrong content type: got %s, want application/grpc", contentType[0])
			}

			if !params.insecure {
				authority, ok := md[":authority"]
				if !ok || len(authority) == 0 {
					return nil, status.Error(codes.DataLoss, "Failed to get authority from metadata")
				}
				if authority[0] != params.authority {
					return nil, status.Errorf(codes.Unauthenticated, "Wrong authority: %s, want %s", authority[0], params.authority)
				}

				wantToken := params.token.TokenType + " " + params.token.AccessToken

				token, ok := md["authorization"]
				if !ok || len(token) == 0 {
					return nil, status.Error(codes.DataLoss, "Failed to get token from metadata")
				}
				if token[0] != wantToken {
					return nil, status.Errorf(codes.Unauthenticated, "Wrong token: %s, want %s", token[0], wantToken)
				}
			}

			return new(testgrpc.Empty), nil
		},
	}
}

func TestConnDevOK(t *testing.T) {
	arpc.SystemCertPool = arpcmocks.ClientCerts()
	arpc.NewTokenSource = arpcmocks.TokenSource(nil)

	stubbedServer := setupClientStubServer(t, stubServerParams{insecure: true})
	clean, err := arpcmocks.Server(stubbedServer, nil, nil)
	require.NoError(t, err)
	defer clean()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	testutils.RequireGRPCCodesEqual(t, err, codes.OK)
}

func TestConnDevErrorSecureServer(t *testing.T) {
	arpc.SystemCertPool = arpcmocks.ClientCerts(x509mocks.ServerCACertPEM)
	arpc.NewTokenSource = arpcmocks.TokenSource(new(arpcmocks.IDTokenStub))

	stubbedServer := setupClientStubServer(
		t,
		stubServerParams{insecure: false, authority: "127.0.0.1:8080", token: arpcmocks.DefaultToken},
	)
	clean, err := arpcmocks.Server(stubbedServer, x509mocks.Server1KeyPEM, x509mocks.Server1CertPEM)
	require.NoError(t, err)
	defer clean()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	require.Error(t, err)
}

func TestConnReleaseOK(t *testing.T) {
	arpc.SystemCertPool = arpcmocks.ClientCerts(x509mocks.ServerCACertPEM)
	arpc.NewTokenSource = arpcmocks.TokenSource(new(arpcmocks.IDTokenStub))

	stubbedServer := setupClientStubServer(
		t,
		stubServerParams{insecure: false, authority: "127.0.0.1:8080", token: arpcmocks.DefaultToken},
	)
	clean, err := arpcmocks.Server(stubbedServer, x509mocks.Server1KeyPEM, x509mocks.Server1CertPEM)
	require.NoError(t, err)
	defer clean()

	connPool := arpc.NewConnPool(true)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	testutils.RequireGRPCCodesEqual(t, err, codes.OK)
}

func TestCloseConn(t *testing.T) {
	arpc.SystemCertPool = arpcmocks.ClientCerts()
	arpc.NewTokenSource = arpcmocks.TokenSource(nil)

	stubbedServer := setupClientStubServer(t, stubServerParams{insecure: true})
	clean, err := arpcmocks.Server(stubbedServer, nil, nil)
	require.NoError(t, err)
	defer clean()

	connPool := arpc.NewConnPool(false)
	defer connPool.Close()

	conn, err := connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.NoError(t, err)

	client := testgrpc.NewTestServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	testutils.RequireGRPCCodesEqual(t, err, codes.OK)

	connPool.Close()

	_, err = client.EmptyCall(ctx, new(testgrpc.Empty))
	require.Error(t, err)

	// Closing multiple times should not cause any issue.
	connPool.Close()
}

func TestOpenClosedPool(t *testing.T) {
	arpc.SystemCertPool = arpcmocks.ClientCerts(x509mocks.ServerCACertPEM)
	arpc.NewTokenSource = arpcmocks.TokenSource(new(arpcmocks.IDTokenStub))

	stubbedClient := setupClientStubServer(t, stubServerParams{insecure: true})
	clean, err := arpcmocks.Server(stubbedClient, x509mocks.Server1KeyPEM, x509mocks.Server1CertPEM)
	require.NoError(t, err)
	defer clean()

	connPool := arpc.NewConnPool(false)
	connPool.Close()

	_, err = connPool.Open(context.Background(), "127.0.0.1", 8080, arpc.ProtocolHTTPS)
	require.ErrorIs(t, err, arpc.ErrConnectionPoolClosed)
}
