package main

import (
	"log"
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

	"service-booking/config"
	"service-booking/db"
	"service-booking/routes"
)

func main() {
	// Determine the environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load .env file
	loadEnvFile(env)

	// Set Gin mode based on environment
	setGinMode(env)

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Register the MySQL database
	db.RegisterMySQL()

	// Defer database connection closure
	sqlDB, err := db.GetDB().DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// Set up routes and start the server
	router := routes.SetupRouter()

	// Determine port (environment variable takes precedence)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = config.String("httpport")
	}

	// Prepare server address
	addr := ":" + port

	// Log server startup
	log.Printf("Starting %s server on port %s...", env, port)

	// Run the server
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile(env string) {
	// Try to load environment-specific .env file first
	envFile := ".env"
	if env != "development" {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	// Load the .env file
	if err := godotenv.Load(envFile); err != nil {
		// Only log if it's not just that the file doesn't exist
		if !os.IsNotExist(err) {
			log.Printf("Error loading %s file: %v", envFile, err)
		}
		
		// Attempt to load default .env
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}
}

// setGinMode sets the Gin framework mode based on the environment
func setGinMode(env string) {
	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}