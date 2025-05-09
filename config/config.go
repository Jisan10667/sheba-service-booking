package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"gopkg.in/ini.v1"
)

var AppConfig struct {
	MySQL struct {
		ServiceBookingDBConn string
		Host                 string
		Port                 string
		User                 string
		Password             string
		Name                 string
	}
	HttpPort int
	JWT      struct {
		SecretKey            string
		AccessTokenDuration  string
		RefreshTokenDuration string
	}
}

// LoadConfig loads the configuration from the app.conf file located in the config folder
func LoadConfig() error {
	// Get the absolute path to the config folder and file
	configFilePath := filepath.Join("config", "app.conf")

	// Check if the file exists
	_, err := os.Stat(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
		return err
	}

	// Parse the INI configuration file
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
		return err
	}

	// Load MySQL configuration with environment variable fallback
	AppConfig.MySQL.Host = getEnv("DB_HOST", cfg.Section("mysql").Key("db_host").String(), "localhost")
	AppConfig.MySQL.Port = getEnv("DB_PORT", cfg.Section("mysql").Key("db_port").String(), "3306")
	AppConfig.MySQL.User = getEnv("DB_USER", cfg.Section("mysql").Key("db_user").String(), "root")
	AppConfig.MySQL.Password = getEnv("DB_PASSWORD", cfg.Section("mysql").Key("db_password").String(), "")
	AppConfig.MySQL.Name = getEnv("DB_NAME", cfg.Section("mysql").Key("db_name").String(), "sheba_service_booking_db")

	// Generate connection string using environment variables
	AppConfig.MySQL.ServiceBookingDBConn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MySQL.User,
		AppConfig.MySQL.Host,
		AppConfig.MySQL.Port,
		AppConfig.MySQL.Name,
	)

	// Load HTTP port
	httpport, err := cfg.Section("").Key("httpport").Int()
	if err != nil {
		log.Printf("Error parsing httpport: %v, using default", err)
		AppConfig.HttpPort = 8087
	} else {
		AppConfig.HttpPort = httpport
	}

	// Load JWT configuration
	AppConfig.JWT.SecretKey = getEnv("JWT_SECRET_KEY", cfg.Section("").Key("jwt_secret_key").String(), "")
	AppConfig.JWT.AccessTokenDuration = getEnv("ACCESS_TOKEN_DURATION", cfg.Section("").Key("access_token_duration").String(), "24h")
	AppConfig.JWT.RefreshTokenDuration = getEnv("REFRESH_TOKEN_DURATION", cfg.Section("").Key("refresh_token_duration").String(), "168h")

	// Logging for debugging
	log.Printf("MySQL Host: %s", AppConfig.MySQL.Host)
	log.Printf("MySQL Port: %s", AppConfig.MySQL.Port)
	log.Printf("MySQL Database: %s", AppConfig.MySQL.Name)
	log.Printf("Connection String: %s", maskConnectionString(AppConfig.MySQL.ServiceBookingDBConn))

	return nil
}


// String returns the value of a configuration key as a string
func String(key string) string {
	log.Printf("Fetching config for key: %s", key)
	switch key {
	case "httpport":
		if AppConfig.HttpPort == 0 {
			log.Printf("Warning: httpport is 0, using default 8087")
			return "8087"
		}
		return fmt.Sprintf("%d", AppConfig.HttpPort)
	case "jwt_secret_key":
		// Check environment variable first
		if secretKey := os.Getenv("JWT_SECRET_KEY"); secretKey != "" {
			return secretKey
		}
		
		// Fallback to AppConfig
		if AppConfig.JWT.SecretKey != "" {
			return AppConfig.JWT.SecretKey
		}
		
		// Generate a fallback secret key
		return generateFallbackSecretKey()
	case "access_token_duration":
		// Check environment variable first
		if duration := os.Getenv("ACCESS_TOKEN_DURATION"); duration != "" {
			return duration
		}
		
		// Fallback to AppConfig
		if AppConfig.JWT.AccessTokenDuration != "" {
			return AppConfig.JWT.AccessTokenDuration
		}
		
		// Default duration
		return "24h"
	case "refresh_token_duration":
		// Check environment variable first
		if duration := os.Getenv("REFRESH_TOKEN_DURATION"); duration != "" {
			return duration
		}
		
		// Fallback to AppConfig
		if AppConfig.JWT.RefreshTokenDuration != "" {
			return AppConfig.JWT.RefreshTokenDuration
		}
		
		// Default duration
		return "168h"
	default:
		log.Printf("Unknown key '%s'", key)
		return ""
	}
}

// generateFallbackSecretKey creates a fallback secret key
func generateFallbackSecretKey() string {
	log.Println("WARNING: Using generated fallback secret key. Set JWT_SECRET_KEY in production!")
	return fmt.Sprintf("fallback-secret-key-%d", time.Now().UnixNano())
}

// getEnv retrieves the value of the environment variable
func getEnv(envVar, configValue, defaultValue string) string {
	// Check environment variable first
	if val := os.Getenv(envVar); val != "" {
		return val
	}

	// Then check config file value
	if configValue != "" {
		return configValue
	}

	// Finally use default value
	return defaultValue
}

// maskConnectionString masks sensitive information in the connection string
func maskConnectionString(dsn string) string {
	// Simple masking of password
	return "mysql://***:***@..." + strings.Split(dsn, "@")[1]
}

// MySQLConfigString generates the MySQL connection string
func MySQLConfigString() string {
	// For a connection without password
	if AppConfig.MySQL.Password == "" {
		return fmt.Sprintf(
			"%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			AppConfig.MySQL.User,
			AppConfig.MySQL.Host,
			AppConfig.MySQL.Port,
			AppConfig.MySQL.Name,
		)
	}

	// For a connection with password
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.MySQL.User,
		AppConfig.MySQL.Password,
		AppConfig.MySQL.Host,
		AppConfig.MySQL.Port,
		AppConfig.MySQL.Name,
	)
}

// Update other methods as needed...