package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-query-builder/querybuilder/pkg/types"
	"github.com/joho/godotenv"
)

// LoadFromEnv loads database configuration from environment variables
func LoadFromEnv() (types.Config, error) {
	// Try to load .env file (ignore errors if file doesn't exist)
	loadDotEnv()
	
	config := types.Config{}

	// Required fields
	driverStr := getEnv("DB_DRIVER", "")
	if driverStr == "" {
		return config, fmt.Errorf("DB_DRIVER environment variable is required")
	}

	switch strings.ToLower(driverStr) {
	case "mysql":
		config.Driver = types.MySQL
	case "postgres", "postgresql":
		config.Driver = types.PostgreSQL
	default:
		return config, fmt.Errorf("unsupported database driver: %s", driverStr)
	}

	config.Host = getEnv("DB_HOST", "localhost")
	config.Database = getEnv("DB_NAME", "")
	if config.Database == "" {
		return config, fmt.Errorf("DB_NAME environment variable is required")
	}

	config.Username = getEnv("DB_USER", "")
	if config.Username == "" {
		return config, fmt.Errorf("DB_USER environment variable is required")
	}

	config.Password = getEnv("DB_PASSWORD", "")

	// Optional fields with defaults
	portStr := getEnv("DB_PORT", "3306")
	if config.Driver == types.PostgreSQL {
		portStr = getEnv("DB_PORT", "5432")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_PORT value: %s", portStr)
	}
	config.Port = port

	config.SSLMode = getEnv("DB_SSL_MODE", "disable")
	config.Charset = getEnv("DB_CHARSET", "utf8mb4")
	config.Timezone = getEnv("DB_TIMEZONE", "UTC")

	// Connection pool settings
	maxOpenConnsStr := getEnv("DB_MAX_OPEN_CONNS", "25")
	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_MAX_OPEN_CONNS value: %s", maxOpenConnsStr)
	}
	config.MaxOpenConns = maxOpenConns

	maxIdleConnsStr := getEnv("DB_MAX_IDLE_CONNS", "5")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_MAX_IDLE_CONNS value: %s", maxIdleConnsStr)
	}
	config.MaxIdleConns = maxIdleConns

	maxLifetimeStr := getEnv("DB_MAX_LIFETIME", "5m")
	maxLifetime, err := time.ParseDuration(maxLifetimeStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_MAX_LIFETIME value: %s", maxLifetimeStr)
	}
	config.ConnMaxLifetime = maxLifetime

	maxIdleTimeStr := getEnv("DB_MAX_IDLE_TIME", "2m")
	maxIdleTime, err := time.ParseDuration(maxIdleTimeStr)
	if err != nil {
		return config, fmt.Errorf("invalid DB_MAX_IDLE_TIME value: %s", maxIdleTimeStr)
	}
	config.ConnMaxIdleTime = maxIdleTime

	return config, nil
}

// LoadFromEnvWithDefaults loads configuration from environment with custom defaults
func LoadFromEnvWithDefaults(defaults types.Config) (types.Config, error) {
	config, err := LoadFromEnv()
	if err != nil {
		return config, err
	}

	// Apply defaults for empty values
	if config.Host == "" && defaults.Host != "" {
		config.Host = defaults.Host
	}
	if config.Port == 0 && defaults.Port != 0 {
		config.Port = defaults.Port
	}
	if config.SSLMode == "" && defaults.SSLMode != "" {
		config.SSLMode = defaults.SSLMode
	}
	if config.Charset == "" && defaults.Charset != "" {
		config.Charset = defaults.Charset
	}
	if config.Timezone == "" && defaults.Timezone != "" {
		config.Timezone = defaults.Timezone
	}
	if config.MaxOpenConns == 0 && defaults.MaxOpenConns != 0 {
		config.MaxOpenConns = defaults.MaxOpenConns
	}
	if config.MaxIdleConns == 0 && defaults.MaxIdleConns != 0 {
		config.MaxIdleConns = defaults.MaxIdleConns
	}
	if config.ConnMaxLifetime == 0 && defaults.ConnMaxLifetime != 0 {
		config.ConnMaxLifetime = defaults.ConnMaxLifetime
	}
	if config.ConnMaxIdleTime == 0 && defaults.ConnMaxIdleTime != 0 {
		config.ConnMaxIdleTime = defaults.ConnMaxIdleTime
	}

	return config, nil
}

// ValidateConfig validates the configuration
func ValidateConfig(config types.Config) error {
	if config.Driver != types.MySQL && config.Driver != types.PostgreSQL {
		return fmt.Errorf("invalid driver: must be mysql or postgresql")
	}

	if config.Host == "" {
		return fmt.Errorf("host is required")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if config.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if config.Username == "" {
		return fmt.Errorf("username is required")
	}

	if config.MaxOpenConns < 0 {
		return fmt.Errorf("max open connections must be non-negative")
	}

	if config.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections must be non-negative")
	}

	if config.MaxIdleConns > config.MaxOpenConns && config.MaxOpenConns > 0 {
		return fmt.Errorf("max idle connections cannot exceed max open connections")
	}

	if config.ConnMaxLifetime < 0 {
		return fmt.Errorf("connection max lifetime must be non-negative")
	}

	if config.ConnMaxIdleTime < 0 {
		return fmt.Errorf("connection max idle time must be non-negative")
	}

	return nil
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SetEnvDefaults sets environment variables to default values if not already set
func SetEnvDefaults() {
	envDefaults := map[string]string{
		"DB_DRIVER":           "mysql",
		"DB_HOST":             "localhost",
		"DB_PORT":             "3306",
		"DB_SSL_MODE":         "disable",
		"DB_CHARSET":          "utf8mb4",
		"DB_TIMEZONE":         "UTC",
		"DB_MAX_OPEN_CONNS":   "25",
		"DB_MAX_IDLE_CONNS":   "5",
		"DB_MAX_LIFETIME":     "5m",
		"DB_MAX_IDLE_TIME":    "2m",
	}

	for key, value := range envDefaults {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// PrintConfig prints the configuration (masking sensitive data)
func PrintConfig(config types.Config) {
	fmt.Printf("Database Configuration:\n")
	fmt.Printf("  Driver: %s\n", config.Driver)
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Database: %s\n", config.Database)
	fmt.Printf("  Username: %s\n", config.Username)
	fmt.Printf("  Password: %s\n", maskPassword(config.Password))
	fmt.Printf("  SSL Mode: %s\n", config.SSLMode)
	fmt.Printf("  Charset: %s\n", config.Charset)
	fmt.Printf("  Timezone: %s\n", config.Timezone)
	fmt.Printf("  Max Open Connections: %d\n", config.MaxOpenConns)
	fmt.Printf("  Max Idle Connections: %d\n", config.MaxIdleConns)
	fmt.Printf("  Connection Max Lifetime: %s\n", config.ConnMaxLifetime)
	fmt.Printf("  Connection Max Idle Time: %s\n", config.ConnMaxIdleTime)
}

func maskPassword(password string) string {
	if password == "" {
		return "(empty)"
	}
	if len(password) <= 2 {
		return "**"
	}
	return password[:1] + strings.Repeat("*", len(password)-2) + password[len(password)-1:]
}

// loadDotEnv attempts to load .env file from current directory or parent directories
func loadDotEnv() {
	// Try current directory first
	if err := godotenv.Load(); err == nil {
		return
	}
	
	// Try parent directories up to 3 levels
	for i := 1; i <= 3; i++ {
		envPath := filepath.Join(strings.Repeat("../", i), ".env")
		if err := godotenv.Load(envPath); err == nil {
			return
		}
	}
	
	// Try absolute paths
	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}
	
	for _, path := range possiblePaths {
		absPath, _ := filepath.Abs(path)
		if err := godotenv.Load(absPath); err == nil {
			return
		}
	}
}