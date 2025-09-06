package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the outbox-api service
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Publish  PublishConfig  `json:"publish"`
	Circuit  CircuitConfig  `json:"circuit"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// PublishConfig holds publishing configuration
type PublishConfig struct {
	BatchSize     int    `json:"batch_size"`
	BatchTimeout  string `json:"batch_timeout"`
	RetryAttempts int    `json:"retry_attempts"`
	RetryDelay    string `json:"retry_delay"`
	MaxRetryDelay string `json:"max_retry_delay"`
	WebhookURL    string `json:"webhook_url"`
}

// CircuitConfig holds circuit breaker configuration
type CircuitConfig struct {
	MaxRequests uint32 `json:"max_requests"`
	Interval    string `json:"interval"`
	Timeout     string `json:"timeout"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  "30s",
			WriteTimeout: "30s",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "outbox",
			SSLMode:  "disable",
		},
		Publish: PublishConfig{
			BatchSize:     10,
			BatchTimeout:  "5s",
			RetryAttempts: 3,
			RetryDelay:    "1s",
			MaxRetryDelay: "30s",
			WebhookURL:    "http://localhost:3000/webhook",
		},
		Circuit: CircuitConfig{
			MaxRequests: 5,
			Interval:    "10s",
			Timeout:     "5s",
		},
	}

	// Load from config file if it exists
	if err := loadFromFile(cfg, "config.json"); err != nil {
		// Config file is optional, so we continue with defaults
		fmt.Printf("Warning: Could not load config file: %v\n", err)
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromFile loads configuration from a JSON file
func loadFromFile(cfg *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}

	if webhookURL := os.Getenv("WEBHOOK_URL"); webhookURL != "" {
		cfg.Publish.WebhookURL = webhookURL
	}
	if batchSize := os.Getenv("BATCH_SIZE"); batchSize != "" {
		if bs, err := strconv.Atoi(batchSize); err == nil {
			cfg.Publish.BatchSize = bs
		}
	}
}

// DSN returns the database connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}
