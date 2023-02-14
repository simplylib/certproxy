package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/simplylib/certproxy/server"
	"golang.org/x/exp/slices"
)

func runServer(ctx context.Context) error {
	args := slices.Delete(append([]string{}, os.Args...), 1, 2)
	flagset := flag.NewFlagSet(args[0], flag.ContinueOnError)

	flagset.Usage = func() {
		fmt.Fprintf(flagset.Output(), "Usage: %v server [flags]\nFlags:\n", args[0])
		flagset.PrintDefaults()
	}

	network := flagset.String("network", ":9777", "network to listen on")
	certPrivPEMPath := flagset.String("certpriv", "", "path to PEM encoded TLS key")
	certPubPEMPath := flagset.String("certpub", "", "path to PEM encoded TLS certificate")

	if err := flagset.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return fmt.Errorf("could not parse command line flags (%w)", err)
	}

	var certificate *tls.Certificate
	if *certPrivPEMPath != "" && *certPubPEMPath != "" {
		cert, err := tls.LoadX509KeyPair(filepath.Clean(*certPubPEMPath), filepath.Clean(*certPrivPEMPath))
		if err != nil {
			return fmt.Errorf("could not load x509 key pair (%w)", err)
		}
		certificate = &cert
	}

	errChan := make(chan error)
	go func() {
		s := server.Server{
			Network:     *network,
			Certificate: certificate,
		}
		err := s.Open()
		if err != nil {
			errChan <- fmt.Errorf("error while running gRPC server (%w)", err)
		}
		errChan <- nil
	}()
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
