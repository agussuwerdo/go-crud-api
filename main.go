package main

import (
	"os"

	_ "goroutine/docs" // Import your docs

	"goroutine/handler"
	// Swagger embed files
)

// PaginatedItemsResponse represents a paginated response for items.
// @Description Represents a paginated response for items
// @Example {"page": 1, "limit": 10, "items": [{"id": "607c191e810c19729de860ea", "name": "Sample Item", "price": 100}]}

func main() {
	handler.InitMongoDB()
	r := handler.SetupRouter()

	// Use the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	r.Run(":" + port)
}
