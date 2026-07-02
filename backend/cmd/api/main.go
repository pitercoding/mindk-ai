package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/routes"
)

func main() {
	routes.RegisterRoutes()

	fmt.Println("\n============== Mindk AI ==============")
	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
