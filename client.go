package main

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
	"log"

	"github.com/simplylib/certproxy/protocol"
	"github.com/simplylib/multierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// generateCSR for requested domain, returning a CSR in DER format, private key in PEM format
// and potentionally an error.
func generateCSR(domain string) (csrder []byte, privKeyPEM []byte, err error) {
	pkey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create ecdsa private key using P384 (%w)", err)
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: domain,
		},
		DNSNames: []string{domain},
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

func getCertificate(ctx context.Context, token string, remote string, domain string) (certificate []byte, privkeypem []byte, err error) {
	log.Printf("INFO: Dialing gRPC endpoint (%v)", remote)
	conn, err := grpc.DialContext(
		ctx,
		remote,
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
	)
	if err != nil {
		return fmt.Errorf("could not Dial gRPC server (%w)", err)
	}
	defer func() {
		if err2 := conn.Close(); err != nil {
			err = multierror.Append(err, fmt.Errorf("could not close gRPC connection (%w)", err2))
		}
	}()

	log.Printf("INFO: Connected to gRPC endpoint, creating certificate request")
	client := protocol.NewCertificateServiceClient(conn)

	reply, err := client.Create(ctx, &protocol.CertificateCreateRequest{})
	if err != nil {
		return fmt.Errorf("could not send a CertificateCreateRequest (%w)", err)
	}

	// verify certificate is valid
	if len(reply.CertPem) == 0 {
		return errors.New("certificate from server is empty")
	}

	cert, err := x509.ParseCertificate(reply.CertPem)
	if err != nil {
		return fmt.Errorf("could not parse certificate from server (%w)", err)
	}

	if cert.Subject.CommonName != domain {
		return fmt.Errorf("common name from certificate (%v) not the same as domain (%v)", cert.Subject.CommonName, domain)
	}

	_ = client
	return nil
}
func runClient(ctx context.Context) error {
	return nil
}
