package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonid6372/notification-processor/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := app.Initialize(ctx)

	app.Start(ctx)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	cancel()

	app.Stop()
}
