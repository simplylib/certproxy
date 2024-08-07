package client

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
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
	return nil, nil
}
