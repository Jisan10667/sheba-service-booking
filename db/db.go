package db

import (
	"fmt"
	"log"
	"service-booking/config"
	"sync"
	"time"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
func RegisterMySQL() {
	var err error
	
	// Create a custom logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
		},
	)

	// Configuration for connection
	dsn := config.MySQLConfigString() // Get the connection string from config
	
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	// Set connection pool parameters
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying DB: %v", err)
	}

	// Connection pool configuration
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	fmt.Println("Connected to MySQL database using GORM")

	// Optional: Auto Migrate models
	// Uncomment and add your models
	// err = DB.AutoMigrate(&model.User{}, &model.Service{}, &model.Booking{})
	// if err != nil {
	// 	log.Fatalf("Error auto-migrating database: %v", err)
	// }
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