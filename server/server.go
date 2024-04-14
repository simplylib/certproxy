package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"

	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/simplylib/certproxy/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthStore interface {
}

// Server handles all gRPC functionalities.
type Server struct {
	// Network is the address that the Server will listen on.
	Network string

	// Certificate to use for listening. If nil, disable TLS.
	Certificate *tls.Certificate

	AuthStore AuthStore

	protocol.UnimplementedCertificateServiceServer
	grpcServer *grpc.Server
}

func (s *Server) Create(ctx context.Context, req *protocol.CertificateCreateRequest) (*protocol.CertificateCreateReply, error) {
	return &protocol.CertificateCreateReply{}, nil
}

// Open GRPC Server and begin handling requests.
func (s *Server) Open() error {
	if s.Certificate == nil {
		slog.Warn("starting server with no certificate, communications are not encrypted")

		s.grpcServer = grpc.NewServer()
	} else {
		slog.Warn("starting server with no certificate")

		s.grpcServer = grpc.NewServer(
			grpc.Creds(credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{*s.Certificate},
				ClientAuth:   tls.NoClientCert,
				MinVersion:   tls.VersionTLS13,
			})),
		)
	}

	listen, err := net.Listen("tcp", s.Network)
	if err != nil {
		return fmt.Errorf("could not listen on tcp socket (%w)", err)
	}

	slog.Info("running gRPC", "endpoint", listen.Addr().String())

	protocol.RegisterCertificateServiceServer(s.grpcServer, s)

	err = s.grpcServer.Serve(listen)
	if err != nil {
		return fmt.Errorf("error while serving grpc server (%w)", err)
	}

	return nil
}

// Close server gracefully, making sure to finish current requests before returning.
func (s *Server) Close() {
	s.grpcServer.GracefulStop()
}
