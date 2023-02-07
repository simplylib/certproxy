package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

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
