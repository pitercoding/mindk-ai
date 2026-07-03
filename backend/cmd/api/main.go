package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/database"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
)

func main() {
	// 1. Database Conectivity
	err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	fmt.Println("\n============== Mindk AI ==============")
	fmt.Println("Database connected successfully")
	fmt.Println("Server running on http://localhost:8080")

	// 2. Register Routes
	routes.RegisterRoutes()

	// 3. Start Server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
