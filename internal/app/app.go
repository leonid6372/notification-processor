package app

import (
	"context"
	"flag"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/config"
	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/internal/domains/scheduler"
	"github.com/leonid6372/notification-processor/internal/domains/sender"
	"github.com/leonid6372/notification-processor/internal/integrations/kafka"
	"github.com/leonid6372/notification-processor/internal/repositories/postgres"
	"github.com/leonid6372/notification-processor/pkg/goosemigrate"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type App struct {
	Config            *config.Config
	Consumer          *kafka.Consumer
	NotificationsRepo domains.NotificationsRepo
	Sender            domains.Sender
	Scheduler         domains.Scheduler

	pool *pgxpool.Pool
}

func Initialize(ctx context.Context) *App {
	log.Info("app initializing...")

	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "bot config path")
	flag.Parse()

	cfg := config.GetConfig(configPath)

	log.Info("init postgres...")

	pool, err := pgxpool.New(ctx, cfg.GetPostgresURL())
	if err != nil {
		log.Fatal("postgres init failed", zap.Error(err))
	}

	if err := goosemigrate.NewMigrator(cfg.GetPostgresURL(), "migrations", cfg.Postgres.Schema).
		Up(); err != nil {
		log.Fatal("migrations up failed", zap.Error(err))
	}

	notificationsRepo := postgres.NewNotificationsRepo(pool)

	log.Info("init kafka consumer...")

	consumer, err := kafka.NewKafkaConsumer(
		ctx, []string{fmt.Sprintf("%s:%d", cfg.Postgres.Host, cfg.Postgres.Port)}, cfg.Kafka.GroupID,
	)
	if err != nil {
		log.Fatal("kafka init failed", zap.Error(err))
	}

	log.Info("init sender...")

	sender := sender.NewSender(ctx, cfg, notificationsRepo)

	log.Info("init scheduler...")

	scheduler := scheduler.NewScheduler(cfg, sender, notificationsRepo)

	return &App{
		Config:            cfg,
		NotificationsRepo: notificationsRepo,
		Consumer:          consumer,
		Sender:            sender,
		Scheduler:         scheduler,

		pool: pool,
	}
}

func (a *App) Start(ctx context.Context) {
	go a.Consumer.Start(
		ctx, []string{a.Config.Kafka.Topic}, kafka.NewHandler(a.Config.Kafka.BatchSize, a.NotificationsRepo),
	)

	go a.Scheduler.StartSendings(ctx)
	go a.Scheduler.StartCleaning(ctx)
}

// Need about 60 seconds to gracefull shoutdown
func (a *App) Stop() {
	log.Info("notification-processor shutting down...")

	if err := a.Consumer.Stop(); err != nil {
		log.Error("kafka consumer stop failed", zap.Error(err))
	}

	a.Scheduler.WaitToStopSending()

	a.pool.Close()
}
