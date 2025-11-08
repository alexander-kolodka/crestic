package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexander-kolodka/crestic/cmd"
)

const gracefulShutdownTimeout = time.Second * 1

func main() {
	os.Exit(execute())
}

func execute() int {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Execute(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return 1
		}
		return 0

	case sig := <-sigCh:
		fmt.Fprintf(os.Stderr, "received %s, shutting down gracefully...\n", sig)

		select {
		case err := <-errCh:
			if err != nil {
				fmt.Fprintf(os.Stderr, "command failed during shutdown: %v\n", err)
			}
			return signalExitCode(sig)

		case <-time.After(gracefulShutdownTimeout):
			fmt.Fprintln(os.Stderr, "graceful shutdown timed out")
			return signalExitCode(sig)
		}
	}
}

// return conventional exit code based on given UNIX signal: 128 + signal number.
func signalExitCode(sig os.Signal) int {
	s, ok := sig.(syscall.Signal)
	if !ok {
		s = syscall.SIGINT
	}

	const exitCode = 128
	return exitCode + int(s)
}
