package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL     string
	HTTPPort        string
	NatsClusterID   string
	NatsClientID    string
	NatsURL         string
	NatsSubject     string
	LogLevel        string
	GracefulTimeout time.Duration
}

func Load() *Config {
	c := &Config{
		DatabaseURL:     getenv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/wb_labs_l0?sslmode=disable"),
		HTTPPort:        getenv("HTTP_PORT", ":8080"),
		NatsClusterID:   getenv("STAN_CLUSTER_ID", "test-cluster"),
		NatsClientID:    getenv("STAN_CLIENT_ID", "wb_labs_l0_backend-"+randSuffix()),
		NatsURL:         getenv("STAN_URL", "nats://nats-streaming:4222"),
		NatsSubject:     getenv("STAN_SUBJECT", "orders"),
		LogLevel:        getenv("LOG_LEVEL", "info"),
		GracefulTimeout: 15 * time.Second,
	}
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func randSuffix() string {
	return time.Now().Format("20060102-150405")
}
