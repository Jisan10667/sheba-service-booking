package main

import (
	"log"
	"service-booking/config"
	"service-booking/db"
	"service-booking/routes"
)

func main() {
	// Load configuration from app.conf
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Register the MySQL database using the configuration
	db.RegisterMySQL()

	// Set up routes and start the server
	router := routes.SetupRouter()

	// Fetch the httpport from the config and run the server
	port := config.String("httpport") // Get the port from config
	log.Printf("Starting server on port %s...", port)
	log.Fatal(router.Run(":" + port)) // Start the server with the correct port
}