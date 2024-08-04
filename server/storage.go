package server

import (
	"context"
)

type Requestor struct {
	ID     string
	ApiKey string
}

type Certificate struct {
	// ID of the Certificate when talking to certproxy server
	ID string `json:"id"`
	// Data of the Certificate in PEM encoding
	PEM []byte `json:"pem"`
}

type RenewError struct {
	NotTime bool
}

type Storage interface {
	// IssueCertificate from Requestor with the domains listed, with idempotency key idemKey.
	IssueCertificate(ctx context.Context, requestor Requestor, domains []string, idemKey string) (Certificate, error)
	// RenewCertificate ID from Requestor with idempotency key idemKey
	RenewCertificate(ctx context.Context, requestor Requestor, id string, idemKey string) (Certificate, error)
	// TODO: RevokeCertificate
}
