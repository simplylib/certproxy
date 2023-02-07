package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/simplylib/certproxy/protocol"
	"github.com/simplylib/multierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func getCertificate(ctx context.Context, remote string, domain string) (err error) {
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

func run() error {
	log.SetFlags(log.Ldate | log.Ltime)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGTERM, os.Interrupt)

		s := <-osSignal
		log.Printf("Cancelling operations due to (%v)\n", s.String())
		cancelFunc()
	}()

	if len(os.Args) < 2 {
		return fmt.Errorf("need subcommand run %v with -h", os.Args[0])
	}

	switch os.Args[1] {
	case "server":
		return runServer(ctx)
	case "client":
		return runClient(ctx)
	default:
	}

	log.SetFlags(0)
	log.Print(
		os.Args[0]+" runs a server or cli client for a certificate proxy\n",
		"\nUsage: "+os.Args[0]+" [command] [flags]\n",
		"\nCommands: server, client\n",
	)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatal(err)
	}
}
