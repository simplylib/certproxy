package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/simplylib/certproxy/client"
)

func run() error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		slog.Info("Cancelling due to interrupt")
	}()

	helpMessage := fmt.Sprintf("%v runs a server or cli client for a certificate proxy\n\nUsage: %v [command] [flags]\n\nCommands: server, client", os.Args[0], os.Args[0])

	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		return errors.New("")
	}

	switch os.Args[1] {
	case "server":
		return runServer(ctx)
	case "client":
		return client.Run(ctx)
	default:
	}

	fmt.Println(helpMessage)
	return errors.New("")
}

func main() {
	if err := run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
