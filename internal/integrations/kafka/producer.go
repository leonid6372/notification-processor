package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/leonid6372/notification-processor/pkg/errs"
	"github.com/leonid6372/notification-processor/pkg/log"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func NewKafkaAsyncProducer(ctx context.Context, brokers []string) (*AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errs.NewStack(err)
	}

	return &AsyncProducer{producer: producer}, nil
}

func (p *AsyncProducer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case success := <-p.producer.Successes():
			log.Debug(
				"message sent", zap.Int32("partition", success.Partition), zap.Int64("offset", success.Offset),
			)

		case err := <-p.producer.Errors():
			log.Error("message sending error", zap.Error(err))
		}
	}
}

func (p *AsyncProducer) Stop() error {
	if err := p.producer.Close(); err != nil {
		return errs.NewStack(err)
	}

	return nil
}

func (p *AsyncProducer) AsyncProduce(msg *sarama.ProducerMessage) {
	p.producer.Input() <- msg
}
