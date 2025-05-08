package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var AppConfig struct {
	MySQL struct {
		ServiceBookingDBConn string
	}
	HttpPort int
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

	// Debugging: Log the value of HttpPort to ensure it's correctly loaded
	log.Printf("Loaded httpport: %d", AppConfig.HttpPort)

	return nil
}

// String returns the value of a configuration key as a string
func String(key string) string {
	log.Printf("Fetching config for key: %s", key) // Added logging for debugging
	switch key {
	case "httpport":
		if AppConfig.HttpPort == 0 {
			log.Printf("Warning: httpport is 0, using default 8087")
			return "8087"
		}
		return fmt.Sprintf("%d", AppConfig.HttpPort)
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