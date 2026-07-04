package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/database"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
)

func main() {
	// 1. Database
	err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. Migrations
	err = migrations.Run(database.DB)
	if err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	fmt.Println("\n============== Mindk AI ==============")
	fmt.Println("Database connected successfully")
	fmt.Println("Server running on http://localhost:8080")

	// 3. Routes
	routes.RegisterRoutes()

	// 4. Server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
