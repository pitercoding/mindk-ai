package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// defaultFrontendOrigin is used when FRONTEND_ORIGIN is unset, so local
// development keeps working without every developer having to set it.
const defaultFrontendOrigin = "http://localhost:5173"

type Config struct {
	OpenAIAPIKey   string
	ClerkSecretKey string
	FrontendOrigin string
}

func Load() *Config {
	err := godotenv.Load() // Search and open .env
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = defaultFrontendOrigin
	}

	cfg := &Config{
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		ClerkSecretKey: os.Getenv("CLERK_SECRET_KEY"),
		FrontendOrigin: frontendOrigin,
	}

	if cfg.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	if cfg.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is not set")
	}

	return cfg
}
