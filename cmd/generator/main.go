package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/internal/config"
	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/internal/integrations/kafka"
	"github.com/leonid6372/notification-processor/pkg/log"
)

const (
	kafkaTopic         = "notifications"
	notificationsCount = 10000
	shippingAddress    = "street Pushkina, building Kolotushkina"
)

var (
	notificationsTypes = []string{
		domains.NotificationTypeOrderCreated,
		domains.NotificationTypePaymentReceived,
		domains.NotificationTypeOrderShipped,
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("generator initializing...")

	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "generator config path")
	flag.Parse()

	cfg := config.GetConfig(configPath)

	log.Info("init kafka producer...")

	producer, err := kafka.NewKafkaAsyncProducer(
		ctx, []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
	)
	if err != nil {
		log.Fatal("kafka producer init failed", zap.Error(err))
	}

	go producer.Start(ctx)

	generatedStat := generateAndProduce(notificationsCount, producer)

	log.Info("generating finished", zap.Any("generatedStat", generatedStat))

	if err := producer.Stop(); err != nil {
		log.Fatal("kafka producer stop failed", zap.Error(err))
	}

	log.Info("producing finished")
}

type Stat struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Total      int
	Correct    int
	BrokenJSON int
}

func generateAndProduce(count int, producer *kafka.AsyncProducer) *Stat {
	stat := Stat{
		StartedAt: time.Now(),
	}

	for range count {
		eventType := notificationsTypes[rand.IntN(3)]
		payload := domains.Payload{}

		switch eventType {
		case domains.NotificationTypeOrderCreated:
			payload.OrderID = uuid.New()

		case domains.NotificationTypePaymentReceived:
			payload.PaymentID = uuid.New()

		case domains.NotificationTypeOrderShipped:
			payload.OrderID = uuid.New()
			payload.ShippingAddress = shippingAddress
		}

		rawPayload, err := json.Marshal(&payload)
		if err != nil {
			log.Error("json.Marshal(&payload) error")
		}

		notification := kafka.Notification{
			EventID:   uuid.New(),
			UserID:    rand.IntN(100000),
			EventType: eventType,
			Timestamp: time.Now().Format(time.RFC3339),
			Payload:   rawPayload,
		}

		data, err := json.Marshal(&notification)
		if err != nil {
			log.Error("json.Marshal(&notification) error")
		}

		var value sarama.ByteEncoder
		if rand.Float64() < 0.02 { // 2% chance to send brokenJSON
			value = sarama.ByteEncoder([]byte{43, 23, 25})
			stat.BrokenJSON++
		} else {
			value = sarama.ByteEncoder(data)
			stat.Correct++
		}

		producer.AsyncProduce(&sarama.ProducerMessage{
			Topic: kafkaTopic,
			Value: value,
		})

		stat.Total++
	}

	stat.FinishedAt = time.Now()

	return &stat
}
