package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	App      AppConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	JWTSecret   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := getEnv("DATABASE_URL", "")
	host := getEnv("DB_HOST", getEnv("POSTGRES_HOST", "localhost"))
	port := getEnvAsInt("DB_PORT", getEnvAsInt("POSTGRES_PORT", 5432))
	user := getEnv("DB_USER", getEnv("POSTGRES_USER", "postgres"))
	password := getEnv("DB_PASSWORD", getEnv("POSTGRES_PASSWORD", ""))
	dbname := getEnv("DB_NAME", getEnv("POSTGRES_DB", "core_bank"))
	sslmode := getEnv("DB_SSL_MODE", "disable")

	config := &Config{
		Database: DatabaseConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			DBName:   dbname,
			SSLMode:  sslmode,
			URL:      dbURL,
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			JWTSecret:   getEnv("JWT_SECRET", os.Getenv("JWT_SECRET")),
		},
	}

	return config, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as an integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
