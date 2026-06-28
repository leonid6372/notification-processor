package scheduler

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/config"
	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type Scheduler struct {
	sender            domains.Sender
	notificationsRepo domains.NotificationsRepo
	notificationsCh   chan *domains.Notification
	doneCh            chan struct{}

	notificationsLimit int
	workersCount       int
	taskTimeout        time.Duration
}

func NewScheduler(
	cfg *config.Config, sender domains.Sender, notificationsRepo domains.NotificationsRepo,
) domains.Scheduler {
	return &Scheduler{
		sender:             sender,
		notificationsRepo:  notificationsRepo,
		notificationsCh:    make(chan *domains.Notification),
		notificationsLimit: cfg.Scheduler.BatchSize,
		workersCount:       cfg.Scheduler.WorkersCount,
		taskTimeout:        cfg.Scheduler.TaskTimeout,
	}
}

func (s *Scheduler) StartSendings(ctx context.Context) {
	for range s.workersCount {
		go s.sender.StartSenderWorker(s.notificationsCh, s.doneCh)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("scheduler sending shutting down...")
			close(s.notificationsCh)
			return

		default:
			notifications, err := s.notificationsRepo.GetNotificationsToSend(ctx, s.notificationsLimit)
			if err != nil {
				log.Error("failed to get notifications to send", zap.Error(err))
				continue
			}

			for _, notification := range notifications {
				s.notificationsCh <- notification
			}

			time.Sleep(10 * time.Second)
		}
	}
}

func (s *Scheduler) WaitToStopSending() {
	for range s.workersCount {
		<-s.doneCh
	}
}

func (s *Scheduler) StartCleaning(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info("scheduler cleaning shutting down...")
			return

		default:
			if err := s.notificationsRepo.ResetZombieNotifications(ctx, s.taskTimeout); err != nil {
				log.Error("failed to reset zombie notifications", zap.Error(err))
			}

			time.Sleep(60 * time.Second)
		}
	}
}
