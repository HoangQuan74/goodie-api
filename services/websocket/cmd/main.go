package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/HoangQuan74/goodie-api/pkg/logger"
	"github.com/HoangQuan74/goodie-api/services/websocket/config"
	"go.uber.org/zap"
)

// Hub manages all WebSocket connections
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool // channel -> connections
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *Hub) Subscribe(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[channel] == nil {
		h.clients[channel] = make(map[*websocket.Conn]bool)
	}
	h.clients[channel][conn] = true
}

func (h *Hub) Unsubscribe(channel string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[channel]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, channel)
		}
	}
}

func (h *Hub) Broadcast(channel string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conns, ok := h.clients[channel]; ok {
		for conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Get().Error("failed to send ws message", zap.Error(err))
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: restrict origins in production
	},
}

func main() {
	cfg := config.Load()

	if err := logger.Init(logger.Config{
		Level:       "info",
		ServiceName: "websocket-service",
		Environment: cfg.Env,
	}); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	hub := NewHub()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "websocket-service",
		})
	})

	// WebSocket endpoint
	// Usage: ws://localhost:8084/ws?channel=order:uuid-123
	router.GET("/ws", func(c *gin.Context) {
		channel := c.Query("channel")
		if channel == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "channel is required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("websocket upgrade failed", zap.Error(err))
			return
		}
		defer conn.Close()

		hub.Subscribe(channel, conn)
		defer hub.Unsubscribe(channel, conn)

		log.Info("client connected", zap.String("channel", channel))

		// Read loop (keep connection alive, handle client messages)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Error("websocket read error", zap.Error(err))
				}
				break
			}
		}

		log.Info("client disconnected", zap.String("channel", channel))
	})

	// Internal endpoint for other services to push messages
	router.POST("/internal/broadcast", func(c *gin.Context) {
		var req struct {
			Channel string `json:"channel" binding:"required"`
			Message string `json:"message" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hub.Broadcast(req.Channel, []byte(req.Message))
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("websocket service starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down websocket service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("websocket service stopped")
}
