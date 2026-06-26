package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/config"
	"github.com/leonid6372/notification-processor/internal/repositories/postgres"
	"github.com/leonid6372/notification-processor/pkg/goosemigrate"
	"github.com/leonid6372/notification-processor/pkg/log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "bot config path")
	flag.Parse()

	cfg := config.GetConfig(configPath)

	ctx, cancel := context.WithCancel(context.Background())

	log.Info("notification-processor starting...")

	log.Info("init postgres...")
	pool, err := pgxpool.New(ctx, cfg.GetPostgresURL())
	if err != nil {
		log.Fatal("postgres init failed", zap.Error(err))
	}

	if err := goosemigrate.NewMigrator(cfg.GetPostgresURL(), "migrations", cfg.Postgres.Schema).
		Up(); err != nil {
		log.Fatal("migrations up failed", zap.Error(err))
	}

	notificationsRepository := postgres.NewNotificationsRepository(pool)

	log.Info("init notification-processor...")
	// app := app.NewApp()

	go func() {
		// app.Start()
	}()

	log.Info("notification-processor starting complete")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Info("notification-processor shutting down...")

	pool.Close()
	// app.Stop()

	if err := log.Sync(); err != nil {
		log.Error("log sync failed", zap.Error(err))
	}

	cancel()

	log.Info("notification-processor shut down complete")
}
