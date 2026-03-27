package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (c Config) URI() string {
	if c.User != "" && c.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
	}
	return fmt.Sprintf("mongodb://%s:%d", c.Host, c.Port)
}

func NewClient(ctx context.Context, cfg Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI())
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("connect mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("ping mongodb: %w", err)
	}

	db := client.Database(cfg.Database)

	logger.Get().Info("connected to MongoDB",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	return client, db, nil
}
