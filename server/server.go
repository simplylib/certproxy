package server

import (
	"context"
	"crypto/tls"
	"errors"
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

// Open http Server and begin handling requests. This is blocking.
func (s *Server) Open() error {
	mux := http.ServeMux{}

	s.httpServ = http.Server{
		Addr:              s.Network,
		Handler:           &mux,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Second * 10,
		MaxHeaderBytes:    1024 * 10,
	}

	// TODO: automatically get our own certificate from
	if s.Certificate == nil {
		slog.Warn("starting server with no certificate, communications are in plaintext.", "addr", s.Network)
	} else {
		s.httpServ.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*s.Certificate},
			MinVersion:   tls.VersionTLS13,
		}
		slog.Info("starting server with certificate", "addr", s.Network)
	}

	mux.HandleFunc("POST /issue", s.Issue)

	if err := s.httpServ.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

// TODO: implement replay attack protection
func (s *Server) Issue(w http.ResponseWriter, r *http.Request) {

}
