package config

import "os"

type Config struct {
	Port  string
	Env   string
	Redis RedisConfig
	Kafka KafkaConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type KafkaConfig struct {
	Brokers []string
}

func Load() *Config {
	return &Config{
		Port: getEnv("WEBSOCKET_PORT", "8084"),
		Env:  getEnv("APP_ENV", "development"),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9094")},
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
