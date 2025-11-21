package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	_ "github.com/tarantool/go-tarantool/uuid"
	_ "go.uber.org/automaxprocs"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/cmd/internal"
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
