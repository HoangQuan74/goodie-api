package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env      string
	Postgres PostgresConfig
	Redis    RedisConfig
	Mongo    MongoConfig
	Kafka    KafkaConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type MongoConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type KafkaConfig struct {
	Brokers []string
}

func Load() *Config {
	return &Config{
		Env: getEnv("APP_ENV", "development"),
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "goodie"),
			Password: getEnv("POSTGRES_PASSWORD", "goodie_secret"),
			Database: getEnv("POSTGRES_DB", "goodie"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Mongo: MongoConfig{
			Host:     getEnv("MONGO_HOST", "localhost"),
			Port:     getEnvInt("MONGO_PORT", 27017),
			User:     getEnv("MONGO_USER", "goodie"),
			Password: getEnv("MONGO_PASSWORD", "goodie_secret"),
			Database: getEnv("MONGO_DB", "goodie"),
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

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
