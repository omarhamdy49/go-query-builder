package config

import (
	"os"
	"testing"
	"time"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

func TestLoadFromEnv(t *testing.T) {
	// Set test environment variables
	testEnv := map[string]string{
		"DB_DRIVER":           "mysql",
		"DB_HOST":             "testhost",
		"DB_PORT":             "3308",
		"DB_NAME":             "testdb",
		"DB_USER":             "testuser",
		"DB_PASSWORD":         "testpass",
		"DB_SSL_MODE":         "disable",
		"DB_CHARSET":          "utf8mb4",
		"DB_TIMEZONE":         "UTC",
		"DB_MAX_OPEN_CONNS":   "50",
		"DB_MAX_IDLE_CONNS":   "10",
		"DB_MAX_LIFETIME":     "10m",
		"DB_MAX_IDLE_TIME":    "5m",
	}

	// Save original env vars
	originalEnv := make(map[string]string)
	for key := range testEnv {
		originalEnv[key] = os.Getenv(key)
	}

	// Set test env vars
	for key, value := range testEnv {
		os.Setenv(key, value)
	}

	// Restore original env vars after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	config, err := loadFromEnvInternal(false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify configuration
	if config.Driver != types.MySQL {
		t.Errorf("Expected driver MySQL, got: %s", config.Driver)
	}
	if config.Host != "testhost" {
		t.Errorf("Expected host 'testhost', got: %s", config.Host)
	}
	if config.Port != 3308 {
		t.Errorf("Expected port 3308, got: %d", config.Port)
	}
	if config.Database != "testdb" {
		t.Errorf("Expected database 'testdb', got: %s", config.Database)
	}
	if config.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", config.Username)
	}
	if config.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got: %s", config.Password)
	}
	if config.MaxOpenConns != 50 {
		t.Errorf("Expected max open conns 50, got: %d", config.MaxOpenConns)
	}
	if config.MaxIdleConns != 10 {
		t.Errorf("Expected max idle conns 10, got: %d", config.MaxIdleConns)
	}
	if config.ConnMaxLifetime != 10*time.Minute {
		t.Errorf("Expected conn max lifetime 10m, got: %s", config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime != 5*time.Minute {
		t.Errorf("Expected conn max idle time 5m, got: %s", config.ConnMaxIdleTime)
	}
}

func TestLoadFromEnvPostgreSQL(t *testing.T) {
	// Set test environment variables for PostgreSQL
	testEnv := map[string]string{
		"DB_DRIVER":   "postgresql",
		"DB_HOST":     "localhost",
		"DB_NAME":     "testdb",
		"DB_USER":     "testuser",
		"DB_PASSWORD": "testpass",
	}

	// Save and clear all env vars that could affect the test
	allEnvKeys := []string{
		"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"DB_SSL_MODE", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", 
		"DB_MAX_LIFETIME", "DB_MAX_IDLE_TIME",
	}
	originalEnv := make(map[string]string)
	for _, key := range allEnvKeys {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	
	// Set test env vars
	for key, value := range testEnv {
		os.Setenv(key, value)
	}

	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	config, err := loadFromEnvInternal(false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Driver != types.PostgreSQL {
		t.Errorf("Expected driver PostgreSQL, got: %s", config.Driver)
	}

	// PostgreSQL should default to port 5432
	if config.Port != 5432 {
		t.Errorf("Expected port 5432, got: %d", config.Port)
	}
}

func TestLoadFromEnvMissingRequired(t *testing.T) {
	// Clear all relevant env vars
	envVars := []string{
		"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"DB_SSL_MODE", "DB_CHARSET", "DB_TIMEZONE", "DB_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS", "DB_MAX_LIFETIME", "DB_MAX_IDLE_TIME",
	}

	originalEnv := make(map[string]string)
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}

	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "Missing driver",
			envVars:  map[string]string{},
			expected: "DB_DRIVER environment variable is required",
		},
		{
			name:     "Missing database name",
			envVars:  map[string]string{"DB_DRIVER": "mysql"},
			expected: "DB_NAME environment variable is required",
		},
		{
			name: "Missing username",
			envVars: map[string]string{
				"DB_DRIVER": "mysql",
				"DB_NAME":   "testdb",
			},
			expected: "DB_USER environment variable is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			allEnvKeys := []string{
				"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
				"DB_SSL_MODE", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", 
				"DB_MAX_LIFETIME", "DB_MAX_IDLE_TIME",
			}
			for _, key := range allEnvKeys {
				os.Unsetenv(key)
			}
			defer func() {
				for _, key := range allEnvKeys {
					os.Unsetenv(key)
				}
			}()
			
			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			_, err := loadFromEnvInternal(false)
			if err == nil {
				t.Errorf("Expected error, got nil")
				return
			}

			if err.Error() != tt.expected {
				t.Errorf("Expected error '%s', got '%s'", tt.expected, err.Error())
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    types.Config
		shouldErr bool
		errMsg    string
	}{
		{
			name: "Valid MySQL config",
			config: types.Config{
				Driver:          types.MySQL,
				Host:            "localhost",
				Port:            3306,
				Database:        "testdb",
				Username:        "user",
				Password:        "pass",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 2 * time.Minute,
			},
			shouldErr: false,
		},
		{
			name: "Valid PostgreSQL config",
			config: types.Config{
				Driver:          types.PostgreSQL,
				Host:            "localhost",
				Port:            5432,
				Database:        "testdb",
				Username:        "user",
				MaxOpenConns:    25,
				MaxIdleConns:    5,
				ConnMaxLifetime: 5 * time.Minute,
				ConnMaxIdleTime: 2 * time.Minute,
			},
			shouldErr: false,
		},
		{
			name: "Invalid driver",
			config: types.Config{
				Driver: types.Driver("invalid"),
			},
			shouldErr: true,
			errMsg:    "invalid driver: must be mysql or postgresql",
		},
		{
			name: "Missing host",
			config: types.Config{
				Driver: types.MySQL,
				Port:   3306,
			},
			shouldErr: true,
			errMsg:    "host is required",
		},
		{
			name: "Invalid port",
			config: types.Config{
				Driver: types.MySQL,
				Host:   "localhost",
				Port:   0,
			},
			shouldErr: true,
			errMsg:    "port must be between 1 and 65535",
		},
		{
			name: "Missing database",
			config: types.Config{
				Driver: types.MySQL,
				Host:   "localhost",
				Port:   3306,
			},
			shouldErr: true,
			errMsg:    "database name is required",
		},
		{
			name: "Missing username",
			config: types.Config{
				Driver:   types.MySQL,
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
			},
			shouldErr: true,
			errMsg:    "username is required",
		},
		{
			name: "Invalid max connections",
			config: types.Config{
				Driver:       types.MySQL,
				Host:         "localhost",
				Port:         3306,
				Database:     "testdb",
				Username:     "user",
				MaxOpenConns: 10,
				MaxIdleConns: 20,
			},
			shouldErr: true,
			errMsg:    "max idle connections cannot exceed max open connections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tt.shouldErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Expected error '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

func TestMaskPassword(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "(empty)"},
		{"a", "**"},
		{"ab", "**"},
		{"abc", "a*c"},
		{"password", "p******d"},
		{"verylongpassword", "v**************d"},
	}

	for _, tt := range tests {
		result := maskPassword(tt.input)
		if result != tt.expected {
			t.Errorf("maskPassword(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}