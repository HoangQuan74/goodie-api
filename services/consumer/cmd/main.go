package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kainguyen/goodie-api/pkg/kafka"
	"github.com/kainguyen/goodie-api/pkg/logger"
	"github.com/kainguyen/goodie-api/pkg/mongo"
	"github.com/kainguyen/goodie-api/pkg/postgres"
	pkgredis "github.com/kainguyen/goodie-api/pkg/redis"
	"github.com/kainguyen/goodie-api/services/consumer/config"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	if err := logger.Init(logger.Config{
		Level:       "info",
		ServiceName: "consumer-service",
		Environment: cfg.Env,
	}); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect PostgreSQL
	pgPool, err := postgres.NewPool(ctx, postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
	})
	if err != nil {
		log.Fatal("failed to connect postgres", zap.Error(err))
	}
	defer pgPool.Close()

	// Connect Redis
	_, err = pkgredis.NewClient(ctx, pkgredis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		log.Fatal("failed to connect redis", zap.Error(err))
	}

	// Connect MongoDB
	mongoClient, _, err := mongo.NewClient(ctx, mongo.Config{
		Host:     cfg.Mongo.Host,
		Port:     cfg.Mongo.Port,
		User:     cfg.Mongo.User,
		Password: cfg.Mongo.Password,
		Database: cfg.Mongo.Database,
	})
	if err != nil {
		log.Fatal("failed to connect mongodb", zap.Error(err))
	}
	defer mongoClient.Disconnect(ctx)

	// Define consumers
	consumers := []struct {
		topic   string
		groupID string
		handler kafka.MessageHandler
	}{
		{
			topic:   "order.created",
			groupID: "order-processor",
			handler: func(ctx context.Context, msg kafkago.Message) error {
				log.Info("processing order.created",
					zap.String("key", string(msg.Key)),
					zap.Int("partition", msg.Partition),
				)
				// TODO: validate order, calculate fees, notify merchant
				return nil
			},
		},
		{
			topic:   "notification.send",
			groupID: "notification-sender",
			handler: func(ctx context.Context, msg kafkago.Message) error {
				var payload map[string]interface{}
				if err := json.Unmarshal(msg.Value, &payload); err != nil {
					return err
				}
				log.Info("sending notification",
					zap.Any("payload", payload),
				)
				// TODO: send push notification / SMS / email
				return nil
			},
		},
		{
			topic:   "order.confirmed",
			groupID: "driver-matcher",
			handler: func(ctx context.Context, msg kafkago.Message) error {
				log.Info("matching driver for order",
					zap.String("key", string(msg.Key)),
				)
				// TODO: find nearest driver via Redis Geo, assign to order
				return nil
			},
		},
		{
			topic:   "review.created",
			groupID: "review-processor",
			handler: func(ctx context.Context, msg kafkago.Message) error {
				log.Info("processing review",
					zap.String("key", string(msg.Key)),
				)
				// TODO: recalculate average rating
				return nil
			},
		},
	}

	// Start all consumers
	var wg sync.WaitGroup
	for _, c := range consumers {
		consumer := kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: cfg.Kafka.Brokers,
			GroupID: c.groupID,
			Topic:   c.topic,
		}, c.handler)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := consumer.Start(ctx); err != nil {
				log.Error("consumer stopped with error", zap.Error(err))
			}
		}()

		log.Info("consumer started", zap.String("topic", c.topic), zap.String("group", c.groupID))
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down consumer service...")
	cancel()
	wg.Wait()
	log.Info("consumer service stopped")
}
