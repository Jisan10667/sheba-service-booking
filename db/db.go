package db

import (
	"fmt"
	"log"
	"service-booking/config"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	once   sync.Once
	dbLock sync.Mutex
)

// GetDB returns the GORM database connection
func GetDB() *gorm.DB {
	dbLock.Lock()
	defer dbLock.Unlock()

	once.Do(func() {
		RegisterMySQL()
	})

	return DB
}

// RegisterMySQL registers a MySQL database connection using GORM
func 	RegisterMySQL() {
	var err error
	dsn := config.MySQLConfigString() // Get the connection string from config
	
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	// Set connection pool parameters
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying DB: %v", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)
	
	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)
	
	fmt.Println("Connected to MySQL database using GORM")
}

// CloseDB closes the MySQL database connection
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting underlying DB:", err)
		return
	}
	
	if err := sqlDB.Close(); err != nil {
		log.Println("Error closing DB connection:", err)
	} else {
		fmt.Println("MySQL database connection closed")
	}
}