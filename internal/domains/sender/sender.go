package sender

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/config"
	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type Sender struct {
	ctx               context.Context
	retryCount        int
	minDelay          time.Duration
	notificationsRepo domains.NotificationsRepo
	r                 *rand.Rand // for error emulating
	// Email integration
	// SMS integration
	// Push integration
}

func NewSender(
	ctx context.Context, cfg *config.Config, notificationsRepo domains.NotificationsRepo,
) domains.Sender {
	return &Sender{
		ctx:               ctx,
		retryCount:        cfg.Sender.RetryCount,
		minDelay:          cfg.Sender.MinDelay,
		notificationsRepo: notificationsRepo,
		r:                 rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Sender) StartSenderWorker(input <-chan *domains.Notification, done chan<- struct{}) {
	for notification := range input {
		title, text, err := notification.GetTitleAndText()
		if err != nil {
			log.Error(
				"failed to get notification title and text",
				zap.String("id", notification.ID.String()),
				zap.Error(err),
			)

			if err := s.notificationsRepo.UpdateNotificationStatus(
				s.ctx, notification.ID, domains.NotificationStatusFailed,
			); err != nil {
				log.Error(
					"failed to update notification status to failed",
					zap.String("id", notification.ID.String()),
					zap.Error(err),
				)
			}

			continue
		}

		switch notification.Type {
		case domains.NotificationTypeOrderCreated:
			emailErr := s.doWithRetry(func() error {
				return s.SendEmail(notification.UserID, title, text)
			})
			if emailErr != nil {
				log.Error(
					"failed to send email notification",
					zap.String("id", notification.ID.String()),
					zap.Error(emailErr),
				)
			}

			pushErr := s.doWithRetry(func() error {
				return s.SendPush(notification.UserID, title, text)
			})
			if pushErr != nil {
				log.Error(
					"failed to send push notification",
					zap.String("id", notification.ID.String()),
					zap.Error(pushErr),
				)
			}

			if emailErr == nil && pushErr == nil {
				if err := s.notificationsRepo.UpdateNotificationStatus(
					s.ctx, notification.ID, domains.NotificationStatusCompleted,
				); err != nil {
					log.Error(
						"failed to update notification status to completed",
						zap.String("id", notification.ID.String()),
						zap.Error(err),
					)
				}
			} else {
				if err := s.notificationsRepo.UpdateNotificationStatus(
					s.ctx, notification.ID, domains.NotificationStatusFailed,
				); err != nil {
					log.Error(
						"failed to update notification status to failed",
						zap.String("id", notification.ID.String()),
						zap.Error(err),
					)
				}
			}

		case domains.NotificationTypePaymentReceived:
			err := s.doWithRetry(func() error {
				return s.SendEmail(notification.UserID, title, text)
			})
			if err != nil {
				log.Error(
					"failed to send email notification",
					zap.String("id", notification.ID.String()),
					zap.Error(err),
				)

				if err := s.notificationsRepo.UpdateNotificationStatus(
					s.ctx, notification.ID, domains.NotificationStatusFailed,
				); err != nil {
					log.Error(
						"failed to update notification status to failed",
						zap.String("id", notification.ID.String()),
						zap.Error(err),
					)
				}
			}

			if err := s.notificationsRepo.UpdateNotificationStatus(
				s.ctx, notification.ID, domains.NotificationStatusCompleted,
			); err != nil {
				log.Error(
					"failed to update notification status to completed",
					zap.String("id", notification.ID.String()),
					zap.Error(err),
				)
			}

		case domains.NotificationTypeOrderShipped:
			pushErr := s.doWithRetry(func() error {
				return s.SendPush(notification.UserID, title, text)
			})
			if pushErr != nil {
				log.Error(
					"failed to send push notification",
					zap.String("id", notification.ID.String()),
					zap.Error(pushErr),
				)
			}

			smsErr := s.doWithRetry(func() error {
				return s.SendSMS(notification.UserID, title, text)
			})
			if smsErr != nil {
				log.Error(
					"failed to send sms notification",
					zap.String("id", notification.ID.String()),
					zap.Error(smsErr),
				)
			}

			if pushErr == nil && smsErr == nil {
				if err := s.notificationsRepo.UpdateNotificationStatus(
					s.ctx, notification.ID, domains.NotificationStatusCompleted,
				); err != nil {
					log.Error(
						"failed to update notification status to completed",
						zap.String("id", notification.ID.String()),
						zap.Error(err),
					)
				}
			} else {
				if err := s.notificationsRepo.UpdateNotificationStatus(
					s.ctx, notification.ID, domains.NotificationStatusFailed,
				); err != nil {
					log.Error(
						"failed to update notification status to failed",
						zap.String("id", notification.ID.String()),
						zap.Error(err),
					)
				}
			}

		default:
			log.Error(
				"invalid notification type",
				zap.String("id", notification.ID.String()),
				zap.String("type", notification.Type),
			)
		}
	}

	done <- struct{}{}
}
