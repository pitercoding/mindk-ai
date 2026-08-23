package main

import (
	"fmt"
	"log"
	"net/http"
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

	// 1. Load environment configuration
	cfg := config.Load()

	// 2. Connect to the database
	err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 3. Run database migrations
	err = migrations.Run(database.DB)
	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}

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

	handler := middleware.SecurityHeaders(corsHandler)

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           handler,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	fmt.Println("\n============== Mindk AI ==============")
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Println("Database: connected")
	fmt.Println("Migrations: applied")
	fmt.Printf("Server listening on %s\n", server.Addr)

	// 6. Start HTTP server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
