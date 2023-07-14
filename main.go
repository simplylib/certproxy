package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/simplylib/certproxy/client"
)

func run() error {
	log.SetFlags(0)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGTERM, os.Interrupt)

		s := <-osSignal
		log.Printf("Cancelling operations due to (%v)\n", s.String())
		cancelFunc()
	}()

	helpMessage := fmt.Sprintf("%v runs a server or cli client for a certificate proxy\n\nUsage: %v [command] [flags]\n\nCommands: server, client\n", os.Args[0], os.Args[0])

	if len(os.Args) < 2 {
		return errors.New(helpMessage)
	}

	switch os.Args[1] {
	case "server":
		return runServer(ctx)
	case "client":
		return client.Run(ctx)
	default:
	}

	log.Print(helpMessage)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.SetOutput(os.Stderr)

		if !errors.Is(err, flag.ErrHelp) {
			log.Fatal(err)
		}
	}
}
