package application

import (
	"os"
	"time"
)

type config struct {
	DatabaseURL     string
	Port            string
	AWSEndpoint     string
	ShutdownTimeout time.Duration
}

func loadConfig() config {
	return config{
		DatabaseURL:     getEnvOrDefault("DATABASE_URL", ""),
		Port:            getEnvOrDefault("PORT", "8080"),
		AWSEndpoint:     getEnvOrDefault("AWS_ENDPOINT", "http://localhost:4566"),
		ShutdownTimeout: time.Second * 10,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
