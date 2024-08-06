// main.go
package main

import (
	"os"

	_ "goroutine/docs" // swagger docs folder

	"goroutine/handler"
	// Swagger embed files
)

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
