package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration for the outbox-api service
type Config struct {
	Server       ServerConfig       `json:"server"`
	Database     DatabaseConfig     `json:"database"`
	Publish      PublishConfig      `json:"publish"`
	Circuit      CircuitConfig      `json:"circuit"`
	FeatureFlags FeatureFlagsConfig `json:"feature_flags"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         string   `json:"port"`
	ReadTimeout  string   `json:"read_timeout"`
	WriteTimeout string   `json:"write_timeout"`
	CORSOrigins  []string `json:"cors_origins"`
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

// FeatureFlagsConfig holds feature flag service configuration
type FeatureFlagsConfig struct {
	BaseURL     string `json:"base_url"`
	Environment string `json:"environment"`
}

// Load loads configuration from .env file and environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  "30s",
			WriteTimeout: "30s",
			CORSOrigins:  []string{"http://localhost:3000", "http://portfolio:3000"},
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
			WebhookURL:    "http://localhost:3000/api/webhook",
		},
		Circuit: CircuitConfig{
			MaxRequests: 5,
			Interval:    "10s",
			Timeout:     "5s",
		},
		FeatureFlags: FeatureFlagsConfig{
			BaseURL:     "http://localhost:4000",
			Environment: "local",
		},
	}

	// Load from .env file if it exists
	if err := loadFromEnvFile(".env"); err != nil {
		// .env file is optional, so we continue with defaults
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Override with environment variables
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromEnvFile loads environment variables from a .env file
func loadFromEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
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

	if flagsURL := os.Getenv("FEATURE_FLAGS_API_URL"); flagsURL != "" {
		cfg.FeatureFlags.BaseURL = flagsURL
	}
	if flagsEnv := os.Getenv("FEATURE_FLAGS_ENV"); flagsEnv != "" {
		cfg.FeatureFlags.Environment = flagsEnv
	}

	if corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); corsOrigins != "" {
		// Parse comma-separated origins
		origins := strings.Split(corsOrigins, ",")
		cfg.Server.CORSOrigins = make([]string, 0, len(origins))
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				cfg.Server.CORSOrigins = append(cfg.Server.CORSOrigins, trimmed)
			}
		}
	}
}

// DSN returns the database connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}
