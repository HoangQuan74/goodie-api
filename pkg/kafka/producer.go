package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kainguyen/goodie-api/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ProducerConfig struct {
	Brokers []string
}

type Producer struct {
	writers map[string]*kafka.Writer
	brokers []string
}

func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		brokers: cfg.Brokers,
	}
}

func (p *Producer) getWriter(topic string) *kafka.Writer {
	if w, ok := p.writers[topic]; ok {
		return w
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
	p.writers[topic] = w
	return w
}

type Event struct {
	Key  string
	Data interface{}
}

func (p *Producer) Publish(ctx context.Context, topic string, event Event) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal event data: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.Key),
		Value: data,
		Time:  time.Now(),
	}

	writer := p.getWriter(topic)
	if err := writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("publish to %s: %w", topic, err)
	}

	logger.Get().Debug("published event",
		zap.String("topic", topic),
		zap.String("key", event.Key),
	)

	return nil
}

func (p *Producer) Close() error {
	for topic, w := range p.writers {
		if err := w.Close(); err != nil {
			logger.Get().Error("failed to close kafka writer",
				zap.String("topic", topic),
				zap.Error(err),
			)
		}
	}
	return nil
}
