package arpcmocks

import (
	"crypto/tls"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
)

func Server(srv testgrpc.TestServiceServer, keyFile, certFile []byte) (func(), error) {
	var sOpts []grpc.ServerOption

	if keyFile == nil || certFile == nil {
		sOpts = append(sOpts, grpc.Creds(insecure.NewCredentials()))
	} else {
		cert, err := tls.X509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		transport := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		sOpts = append(sOpts, grpc.Creds(transport))
	}

	s := grpc.NewServer(sOpts...)

	testgrpc.RegisterTestServiceServer(s, srv)

	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.Serve(lis)
	}()

	stop := func() {
		s.Stop()
		_ = lis.Close()
	}

	return stop, nil
}
