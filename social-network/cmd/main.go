package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	_ "go.uber.org/automaxprocs"

	"github.com/shaelmaar/otus-highload/social-network/cmd/internal"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer cancel()

	cobra.CheckErr(internal.Execute(ctx))
}
