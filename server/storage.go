package server

import (
	"context"
)

type Requestor struct {
	ID     string
	ApiKey string
}

type Certificate struct {
	// ID of the Certificate
	ID string `json:"id"`
	// PEM encoding of the Certificate
	PEM []byte `json:"pem"`
}

// TODO: requests should expire
// TODO: version the database and automigrate
type Storage interface {
	// IssueCertificate from Requestor with the domains listed, with idempotency key idemKey.
	IssueCertificate(ctx context.Context, requestor Requestor, domains []string, idemKey string) (Certificate, error)
	// RenewCertificate ID from Requestor with idempotency key idemKey
	RenewCertificate(ctx context.Context, requestor Requestor, id string, idemKey string) (Certificate, error)
	// TODO: RevokeCertificate
}
