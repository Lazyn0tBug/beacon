package database

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// var logger = Logger::GetLogger()

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

type Database struct {
	DB *gorm.DB
}

var Instance *Database

// Connect establishes a connection to the database based on the provided environment
// Connect 根据环境连接数据库并返回数据库实例
func Connect(env string) (*Database, error) {
	var err error
	var config Config

	// 解析配置文件
	err = parseConfigFile(&config)
	if err != nil {
		return nil, err
	}

	// 根据环境获取配置
	cfg := getConfigByEnv(env, &config)

	// 生成DSN字符串
	dsn, err := generateDSN(cfg)
	if err != nil {
		return nil, err
	}

	// 根据方言连接数据库
	db, err := connectDatabase(dsn, cfg.Dialect)
	if err != nil {
		return nil, err
	}

	database := &Database{DB: db}
	Instance = database
	return database, nil
}

// parseConfigFile 解析TOML配置文件并返回Config实例
func parseConfigFile(config *Config) error {
	_, err := toml.DecodeFile("database/database.toml", config)
	if err != nil {
		log.Fatalf("解析TOML配置文件出错: %v", err)
		return err
	}
	return nil
}

// getConfigByEnv 根据环境获取配置
func getConfigByEnv(env string, config *Config) *DatabaseConfig {
	switch env {
	case "development":
		return &config.Development
	case "test":
		return &config.Test
	case "production":
		return &config.Production
	default:
		return &config.Default
	}
}

// generateDSN 根据配置生成DSN字符串
func generateDSN(cfg *DatabaseConfig) (string, error) {
	dsnTemplate := "host={{ .host }} user={{ .username }} password={{ .password }} dbname={{ .database }} port={{ .port }}"
	dsnTemplate = strings.TrimSuffix(dsnTemplate, "}")

	t := template.Must(template.New("dsn").Parse(dsnTemplate))

	var dsn bytes.Buffer
	err := t.Execute(&dsn, cfg)
	if err != nil {
		panic(err)
	}

	return dsn.String(), nil
}

// connectDatabase 根据方言连接数据库并返回*gorm.DB实例
func connectDatabase(dsn string, dialect string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch dialect {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		log.Fatal("不支持的方言")
	}

	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	return db, nil
}

// MigrateSchemas 可以作为单独的迁移方法
func MigrateSchemas(db *gorm.DB) error {
	// 执行迁移逻辑
	if err := db.AutoMigrate().Error; err != nil {
		return fmt.Errorf("failed to migrate schemas: %v", err)
	}
	return nil
}

func (db *Database) GetDB() *gorm.DB {
	return db.DB
}
