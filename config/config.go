package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

var AppConfig struct {
	MySQL struct {
		ServiceBookingDBConn string
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

	// Load MySQL section
	mysqlSection := cfg.Section("mysql")
	AppConfig.MySQL.ServiceBookingDBConn = mysqlSection.Key("service_booking_db_conn").String()

	// Load HTTP port
	httpport, err := cfg.Section("").Key("httpport").Int()
	if err != nil {
		log.Printf("Error parsing httpport: %v, using default", err)
		AppConfig.HttpPort = 8087 // Default from your conf file
	} else {
		AppConfig.HttpPort = httpport
	}

	// Load JWT configuration
	jwtSection := cfg.Section("")
	AppConfig.JWT.SecretKey = jwtSection.Key("jwt_secret_key").String()
	AppConfig.JWT.AccessTokenDuration = jwtSection.Key("access_token_duration").String()
	AppConfig.JWT.RefreshTokenDuration = jwtSection.Key("refresh_token_duration").String()

	// Debugging: Log the loaded configurations
	log.Printf("Loaded httpport: %d", AppConfig.HttpPort)
	log.Printf("Loaded JWT Secret Key: %s", AppConfig.JWT.SecretKey)
	log.Printf("Loaded Access Token Duration: %s", AppConfig.JWT.AccessTokenDuration)
	log.Printf("Loaded Refresh Token Duration: %s", AppConfig.JWT.RefreshTokenDuration)

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
		return AppConfig.JWT.SecretKey
	case "access_token_duration":
		return AppConfig.JWT.AccessTokenDuration
	case "refresh_token_duration":
		return AppConfig.JWT.RefreshTokenDuration
	default:
		log.Printf("Unknown key '%s'", key)
		return ""
	}
}

// Int returns the value of a configuration key as an integer
func Int(key string) int {
	log.Printf("Fetching int config for key: %s", key)
	switch key {
	case "httpport":
		if AppConfig.HttpPort == 0 {
			log.Printf("Warning: httpport is 0, using default 8087")
			return 8087
		}
		return AppConfig.HttpPort
	default:
		log.Printf("Unknown int key '%s'", key)
		return 0
	}
}

// ParseDuration safely parses a duration string
func ParseDuration(key string) time.Duration {
	log.Printf("Parsing duration for key: %s", key)
	switch key {
	case "access_token_duration":
		duration, err := time.ParseDuration(AppConfig.JWT.AccessTokenDuration)
		if err != nil {
			log.Printf("Error parsing access token duration: %v, using default 24h", err)
			return 24 * time.Hour
		}
		return duration
	case "refresh_token_duration":
		duration, err := time.ParseDuration(AppConfig.JWT.RefreshTokenDuration)
		if err != nil {
			log.Printf("Error parsing refresh token duration: %v, using default 168h", err)
			return 7 * 24 * time.Hour
		}
		return duration
	default:
		log.Printf("Unknown duration key '%s'", key)
		return 0
	}
}

// MySQLConfigString returns the MySQL connection string
func MySQLConfigString() string {
	return AppConfig.MySQL.ServiceBookingDBConn
}

// InitConfig ensures the configuration is loaded before the application starts
func InitConfig() {
	if err := LoadConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
}