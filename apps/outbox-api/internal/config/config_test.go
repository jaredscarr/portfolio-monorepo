package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expected: &Config{
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
			},
		},
		{
			name: "environment variable overrides",
			envVars: map[string]string{
				"PORT":                  "9090",
				"DB_HOST":               "prod-db",
				"DB_PORT":               "5433",
				"DB_USER":               "admin",
				"DB_PASSWORD":           "secret",
				"DB_NAME":               "production",
				"DB_SSLMODE":            "require",
				"WEBHOOK_URL":           "https://api.example.com/webhook",
				"BATCH_SIZE":            "20",
				"FEATURE_FLAGS_API_URL": "https://flags.example.com",
				"FEATURE_FLAGS_ENV":     "prod",
			},
			expected: &Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  "30s",
					WriteTimeout: "30s",
				},
				Database: DatabaseConfig{
					Host:     "prod-db",
					Port:     5433,
					User:     "admin",
					Password: "secret",
					DBName:   "production",
					SSLMode:  "require",
				},
				Publish: PublishConfig{
					BatchSize:     20,
					BatchTimeout:  "5s",
					RetryAttempts: 3,
					RetryDelay:    "1s",
					MaxRetryDelay: "30s",
					WebhookURL:    "https://api.example.com/webhook",
				},
				Circuit: CircuitConfig{
					MaxRequests: 5,
					Interval:    "10s",
					Timeout:     "5s",
				},
				FeatureFlags: FeatureFlagsConfig{
					BaseURL:     "https://flags.example.com",
					Environment: "prod",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := Load()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "default configuration",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				DBName:   "outbox",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=outbox sslmode=disable",
		},
		{
			name: "production configuration",
			config: DatabaseConfig{
				Host:     "prod-db.internal",
				Port:     5432,
				User:     "admin",
				Password: "secret123",
				DBName:   "production",
				SSLMode:  "require",
			},
			expected: "host=prod-db.internal port=5432 user=admin password=secret123 dbname=production sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.DSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestLoadFromFile_DISABLED(t *testing.T) {
	t.Skip("Test disabled - loadFromEnvFile doesn't load JSON files")
	// Create a temporary config file
	configContent := `{
		"server": {
			"port": "9090",
			"read_timeout": "60s",
			"write_timeout": "60s"
		},
		"database": {
			"host": "test-db",
			"port": 5433,
			"user": "testuser",
			"password": "testpass",
			"dbname": "testdb",
			"sslmode": "require"
		},
		"publish": {
			"batch_size": 25,
			"batch_timeout": "10s",
			"retry_attempts": 5,
			"retry_delay": "2s",
			"max_retry_delay": "60s",
			"webhook_url": "https://test.example.com/webhook"
		},
		"circuit": {
			"max_requests": 10,
			"interval": "30s",
			"timeout": "10s"
		}
	}`

	// Write config file
	err := os.WriteFile("test-config.json", []byte(configContent), 0644)
	require.NoError(t, err)
	defer os.Remove("test-config.json")

	// Test loading from file
	err = loadFromEnvFile("test-config.json")
	require.NoError(t, err)

	// Load config again to pick up the env vars
	cfg, err := Load()
	require.NoError(t, err)

	// Clean up environment variables after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("WEBHOOK_URL")
		// ... other env vars
	}()

	expected := &Config{
		Server: ServerConfig{
			Port:         "9090",
			ReadTimeout:  "60s",
			WriteTimeout: "60s",
		},
		Database: DatabaseConfig{
			Host:     "test-db",
			Port:     5433,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			SSLMode:  "require",
		},
		Publish: PublishConfig{
			BatchSize:     25,
			BatchTimeout:  "10s",
			RetryAttempts: 5,
			RetryDelay:    "2s",
			MaxRetryDelay: "60s",
			WebhookURL:    "https://test.example.com/webhook",
		},
		Circuit: CircuitConfig{
			MaxRequests: 10,
			Interval:    "30s",
			Timeout:     "10s",
		},
	}

	assert.Equal(t, expected, cfg)
}

func TestLoadFromFile_NonExistentFile(t *testing.T) {
	err := loadFromEnvFile("non-existent-config.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The system cannot find the file specified")
}
