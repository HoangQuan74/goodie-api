package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	apperrors "github.com/HoangQuan74/goodie-api/pkg/errors"
	"github.com/HoangQuan74/goodie-api/pkg/response"
	"github.com/redis/go-redis/v9"
)

type RateLimitConfig struct {
	Max    int
	Window time.Duration
}

func RateLimit(rdb *redis.Client, cfg RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.FullPath())
		ctx := context.Background()

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, cfg.Window)
		}

		if count > int64(cfg.Max) {
			response.Error(c, apperrors.New(429, "rate limit exceeded"))
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", int64(cfg.Max)-count))
		c.Next()
	}
}
