package kafka

import (
	"context"
	"fmt"

	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

type MessageHandler func(ctx context.Context, msg kafka.Message) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(cfg ConsumerConfig, handler MessageHandler) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	logger.Get().Info("starting consumer",
		zap.String("topic", c.reader.Config().Topic),
		zap.String("group", c.reader.Config().GroupID),
	)

	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				logger.Get().Error("failed to fetch message", zap.Error(err))
				continue
			}

			if err := c.handler(ctx, msg); err != nil {
				logger.Get().Error("failed to handle message",
					zap.String("topic", msg.Topic),
					zap.Int("partition", msg.Partition),
					zap.Int64("offset", msg.Offset),
					zap.Error(err),
				)
				// TODO: send to DLQ
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("commit message: %w", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
