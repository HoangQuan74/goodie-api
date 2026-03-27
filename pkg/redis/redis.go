package redis

import (
	"context"
	"fmt"

	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	logger.Get().Info("connected to Redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return client, nil
}
