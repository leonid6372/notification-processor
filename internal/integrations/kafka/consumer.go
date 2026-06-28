package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/pkg/errs"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

func NewKafkaConsumer(ctx context.Context, brokers []string, groupID string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = false

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	return &Consumer{group: group}, nil
}

func (c *Consumer) Start(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	for {
		if err := c.group.Consume(ctx, topics, handler); err != nil {
			log.Error("kafka consume failed", zap.Error(err))
		}
	}
}

func (c *Consumer) Stop() error {
	if err := c.group.Close(); err != nil {
		return errs.NewStack(err)
	}

	return nil
}
