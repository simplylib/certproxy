package client

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/simplylib/certproxy/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type gRPCConnection struct {
	conn   *grpc.ClientConn
	client protocol.CertificateServiceClient
	token  string
}

func (gc *gRPCConnection) Create(ctx context.Context, csr []byte) ([]byte, error) {
	reply, err := gc.client.Create(ctx, &protocol.CertificateCreateRequest{
		Token:                     gc.token,
		CertificateSigningRequest: csr,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create certificate request with remote server: %w", err)
	}

	return reply.Certificate, nil
}

// TODO: rate limiting
func connectToGRPC(ctx context.Context, host string, token string) (*gRPCConnection, error) {
	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: false,
			}),
		))
	if err != nil {
		return nil, fmt.Errorf("could not grpc.Dial(%v): %w", host, err)
	}

	client := protocol.NewCertificateServiceClient(conn)

	return &gRPCConnection{
		conn:   conn,
		client: client,
		token:  token,
	}, nil
}
