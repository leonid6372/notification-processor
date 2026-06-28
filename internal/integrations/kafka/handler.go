package kafka

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type Handler struct {
	notificationsRepo domains.NotificationsRepo
	eventCount        int
	batchSize         int
}

func NewHandler(batchSize int, notificationsRepo domains.NotificationsRepo) sarama.ConsumerGroupHandler {
	return &Handler{
		notificationsRepo: notificationsRepo,
		batchSize:         batchSize,
	}
}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Info("kafka message channel was closed")
				return nil
			}

			log.Info("kafka consumed event", zap.String("value", string(msg.Value)))

			sess.MarkMessage(msg, "")
			h.eventCount++

			notification := new(Notification)

			err := json.Unmarshal(msg.Value, notification)
			if err != nil {
				log.Error("failed to parse notification from kafka", zap.Error(err))
				continue
			}

			sendAt, err := time.Parse(time.RFC3339, notification.Timestamp)
			if err != nil {
				log.Error(
					"failed to parse notification timestamp from kafka",
					zap.String("id", notification.EventID.String()),
					zap.Error(err),
				)
				continue
			}

			if err := h.notificationsRepo.CreateNotification(sess.Context(), &domains.Notification{
				ID:         notification.EventID,
				UserID:     notification.UserID,
				Type:       notification.EventType,
				RawPayload: notification.Payload,
				SendAt:     sendAt,
			}); err != nil {
				log.Error(
					"failed to create notification in notificationsRepo",
					zap.String("id", notification.EventID.String()),
					zap.Error(err),
				)
				continue
			}

			if h.eventCount >= h.batchSize {
				sess.Commit()
				h.eventCount = 0
			}

		case <-sess.Context().Done():
			return nil
		}
	}
}
