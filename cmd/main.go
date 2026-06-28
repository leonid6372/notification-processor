package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/app"
	"github.com/leonid6372/notification-processor/pkg/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := app.Initialize(ctx)

	app.Start(ctx)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	app.Stop()

	if err := log.Sync(); err != nil {
		log.Error("log sync failed", zap.Error(err))
	}

	cancel()

	log.Info("notification-processor shut down complete")
}
