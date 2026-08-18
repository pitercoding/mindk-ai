package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey   string
	ClerkSecretKey string
}

func Load() *Config {
	err := godotenv.Load() // Search and open .env
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	cfg := &Config{
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		ClerkSecretKey: os.Getenv("CLERK_SECRET_KEY"),
	}

	if cfg.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	if cfg.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is not set")
	}

	return cfg
}
