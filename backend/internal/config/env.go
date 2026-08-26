package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"

	defaultEnvironment = EnvDevelopment

	// defaultFrontendOrigin and defaultDatabasePath are development-only
	// conveniences. In production these must be set explicitly.
	defaultFrontendOrigin = "http://localhost:5173"
	defaultDatabasePath   = "./data/mindk.db"

	// defaultPort is used when the hosting environment (or docker-compose)
	// does not inject PORT. Render sets PORT itself in production.
	defaultPort = "8080"
)

type Config struct {
	Environment    string
	OpenAIAPIKey   string
	ClerkSecretKey string
	FrontendOrigin string

	// Port is the HTTP port the server listens on.
	Port string

	// DatabasePath is the SQLite file used in development. Ignored in
	// production.
	DatabasePath string

	// DatabaseURL is the PostgreSQL connection string used in production.
	// Not used in development.
	DatabaseURL string
}

func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// Load reads a .env file (if present) into the process environment, builds
// the Config, and exits the process if the environment is invalid. A missing
// .env file is not an error: production deployments are expected to inject
// real environment variables instead of shipping a .env file.
func Load() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	cfg, err := FromEnv(os.Getenv)
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	return cfg
}

// FromEnv builds a Config from an arbitrary env lookup function. It applies
// development-only defaults and enforces required/production-only variables,
// but has no process-level side effects (no file I/O, no log.Fatal), so it
// can be exercised directly in tests.
func FromEnv(getenv func(string) string) (*Config, error) {
	environment := getenv("APP_ENV")
	if environment == "" {
		environment = defaultEnvironment
	}
	if environment != EnvDevelopment && environment != EnvProduction {
		return nil, fmt.Errorf("APP_ENV must be %q or %q, got %q", EnvDevelopment, EnvProduction, environment)
	}

	cfg := &Config{
		Environment:    environment,
		OpenAIAPIKey:   getenv("OPENAI_API_KEY"),
		ClerkSecretKey: getenv("CLERK_SECRET_KEY"),
		FrontendOrigin: getenv("FRONTEND_ORIGIN"),
		DatabasePath:   getenv("DATABASE_PATH"),
		DatabaseURL:    getenv("DATABASE_URL"),
		Port:           getenv("PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	if cfg.IsProduction() {
		if cfg.FrontendOrigin == "" {
			return nil, errors.New("FRONTEND_ORIGIN must be set explicitly when APP_ENV=production")
		}
		if cfg.DatabaseURL == "" {
			return nil, errors.New("DATABASE_URL must be set explicitly when APP_ENV=production")
		}
	} else {
		if cfg.FrontendOrigin == "" {
			cfg.FrontendOrigin = defaultFrontendOrigin
		}
		if cfg.DatabasePath == "" {
			cfg.DatabasePath = defaultDatabasePath
		}
	}

	if cfg.OpenAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY is not set")
	}
	if cfg.ClerkSecretKey == "" {
		return nil, errors.New("CLERK_SECRET_KEY is not set")
	}

	return cfg, nil
}
