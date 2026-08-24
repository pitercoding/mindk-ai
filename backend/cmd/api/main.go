package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/database"
	"github.com/pitercoding/mindk-ai/backend/internal/middleware"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
	"github.com/rs/cors"
)

const (
	serverAddr = ":8080"

	// Chat requests wait on a full LLM completion and document uploads can
	// carry up to 10 MiB, so timeouts are generous compared to a typical
	// JSON API, while still bounding how long a stalled connection can hold
	// a server resource open.
	readTimeout       = 30 * time.Second
	readHeaderTimeout = 10 * time.Second
	writeTimeout      = 60 * time.Second
	idleTimeout       = 120 * time.Second
)

func main() {

	// slog.Default is set up first so every subsequent step - including
	// config loading - logs through the same structured JSON handler.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// 1. Load environment configuration
	cfg := config.Load()

	// 2. Connect to the database
	err := database.Connect(cfg.DatabasePath)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected")

	// 3. Run database migrations
	err = migrations.Run(database.DB)
	if err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	// 4. Build application dependencies
	application := app.New(database.DB, cfg)

	// 5. Register HTTP routes
	routes.RegisterRoutes(application)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			cfg.FrontendOrigin,
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	}).Handler(http.DefaultServeMux)

	handler := middleware.RequestLogger(middleware.SecurityHeaders(corsHandler))

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           handler,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	slog.Info("server starting", "environment", cfg.Environment, "addr", server.Addr)

	// 6. Start HTTP server
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
