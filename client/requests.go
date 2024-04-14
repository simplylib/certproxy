package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"

	"github.com/simplylib/certproxy/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// generateCSR for requested domain, returning a CSR in DER format, private key in PEM format
// and potentionally an error.
func generateCSR(domains []string) (csrder []byte, privKeyPEM []byte, err error) {
	pkey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create ecdsa private key using P384 (%w)", err)
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domains[0],
		},
		DNSNames: domains,
	}, pkey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create certificate request (%w)", err)
	}

	pkcs8PrivKey, err := x509.MarshalPKCS8PrivateKey(pkey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal private key to PKCS8 (%w)", err)
	}

	certPrivPEM := &bytes.Buffer{}
	err = pem.Encode(certPrivPEM, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: pkcs8PrivKey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not encode private key pem (%w)", err)
	}

	return csr, certPrivPEM.Bytes(), nil
}

// getCertificate from remote url using token for authentication and csr (in DER format)
func getCertificate(ctx context.Context, token string, remote string, csr []byte) (certificate []byte, err error) {
	slog.Info("gRPC Dialing", "endpoint", remote)
	conn, err := grpc.DialContext(
		ctx,
		remote,
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	)
	if err != nil {
		return nil, fmt.Errorf("could not Dial gRPC server (%w)", err)
	}
	defer func() {
		if err2 := conn.Close(); err != nil {
			err = errors.Join(err, fmt.Errorf("could not close gRPC connection (%w)", err2))
		}
	}()

	slog.Info("Connected to gRPC endpoint, creating certificate request")
	client := protocol.NewCertificateServiceClient(conn)

	reply, err := client.Create(ctx, &protocol.CertificateCreateRequest{
		Token:                     token,
		CertificateSigningRequest: csr,
	})
	if err != nil {
		return nil, fmt.Errorf("could not send a CertificateCreateRequest (%w)", err)
	}

	// TODO: validate certificate
	if len(reply.Certificate) == 0 {
		return nil, errors.New("certificate from server is empty")
	}

	return reply.Certificate, nil
}
