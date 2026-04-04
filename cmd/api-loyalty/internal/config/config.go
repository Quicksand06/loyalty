package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	KafkaBroker string
	KafkaTopic  string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env is optional; real environments inject vars directly
	}

	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return Config{}, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "transactions"
	}

	return Config{
		DatabaseURL: dsn,
		HTTPAddr:    addr,
		KafkaBroker: broker,
		KafkaTopic:  topic,
	}, nil
}
