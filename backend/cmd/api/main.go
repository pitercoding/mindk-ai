package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/database"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
)

func main() {
	// 1. Connect to the database
	err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. Run database migrations
	err = migrations.Run(database.DB)
	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	// 3. Build application dependencies
	application := app.New(database.DB)

	// 4. Register HTTP routes
	routes.RegisterRoutes(application)

	fmt.Println("\n============== Mindk AI ==============")
	fmt.Println("Database connected successfully")
	fmt.Println("Server running on http://localhost:8080")

	// 5. Start HTTP server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
