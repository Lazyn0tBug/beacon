package database

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig represents the database configuration structure
type DatabaseConfig struct {
	Dialect  string
	Database string
	Username string
	Password string
	Host     string
	Port     int
}

// Config holds the parsed TOML configuration
type Config struct {
	Default     DatabaseConfig
	Development DatabaseConfig
	Test        DatabaseConfig
	Production  DatabaseConfig
}

var DB *gorm.DB

// Connect establishes a connection to the database based on the provided environment
func Connect(env string) {
	var config Config
	if _, err := toml.DecodeFile("database/database.toml", &config); err != nil {
		log.Fatalf("Error parsing TOML config file: %v", err)
	}

	var cfg *DatabaseConfig
	switch env {
	case "development":
		cfg = &config.Development
	case "test":
		cfg = &config.Test
	case "production":
		cfg = &config.Production
	default:
		cfg = &config.Default
	}

	var err error
	DB, err = gorm.Open(postgres.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Migrate the schema
	DB.AutoMigrate()
}
