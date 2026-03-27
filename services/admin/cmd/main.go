package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"github.com/HoangQuan74/goodie-api/pkg/middleware"
	"github.com/HoangQuan74/goodie-api/pkg/mongo"
	"github.com/HoangQuan74/goodie-api/pkg/postgres"
	pkgredis "github.com/HoangQuan74/goodie-api/pkg/redis"
	"github.com/HoangQuan74/goodie-api/pkg/validator"
	"github.com/HoangQuan74/goodie-api/services/admin/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Init logger
	if err := logger.Init(logger.Config{
		Level:       "info",
		ServiceName: "admin-service",
		Environment: cfg.Server.Env,
	}); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	ctx := context.Background()

	// Init validator
	validator.Init()

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
	rdb, err := pkgredis.NewClient(ctx, pkgredis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})
	if err != nil {
		log.Fatal("failed to connect redis", zap.Error(err))
	}
	defer rdb.Close()

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

	// Setup Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.Logging(),
	)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "admin-service",
		})
	})

	// API routes
	v1 := router.Group("/api/v1/admin")
	_ = v1 // TODO: register route handlers

	// Start HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("admin service starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down admin service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("admin service stopped")
}
