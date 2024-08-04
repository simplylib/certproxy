package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server handles all gRPC functionalities.
type Server struct {
	// Network is the address that the Server will listen on.
	Network string

	// Certificate to use for listening. If nil, disable TLS.
	Certificate *tls.Certificate

	// Storage to use as backend store.
	Storage Storage

	httpServ http.Server
}

// Open GRPC Server and begin handling requests.
func (s *Server) Open() error {
	mux := http.ServeMux{}

	s.httpServ = http.Server{
		Addr:              s.Network,
		Handler:           &mux,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
		IdleTimeout:       time.Second * 5,
		MaxHeaderBytes:    1024 * 10,
	}

	if s.Certificate == nil {
		slog.Warn("starting server with no certificate, communications are in plaintext.", "addr", s.Network)
	} else {
		s.httpServ.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*s.Certificate},
			MinVersion:   tls.VersionTLS13,
		}
		slog.Info("starting server with certificate", "addr", s.Network)
	}

	mux.Handle("/renew", s)
	mux.Handle("/issue", s)

	if err := s.httpServ.ListenAndServe(); err != nil {
		return fmt.Errorf("httpServ.ListenAndServe error: %w", err)
	}

	return nil
}

// Close server gracefully, making sure to finish current requests before returning.
func (s *Server) Close(ctx context.Context) error {
	if err := s.httpServ.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not httpServ.Shutdown due to net.Listener: %w", err)
	}

	return nil
}
