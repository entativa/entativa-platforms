package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	Security SecurityConfig
	Platform PlatformConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string
	Environment string
	ServiceName string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConnections  int
	MaxIdleConnections int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret              string
	AccessTokenExpiry   time.Duration
	RefreshTokenExpiry  time.Duration
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// SecurityConfig holds security settings
type SecurityConfig struct {
	BcryptCost                int
	PasswordMinLength         int
	PasswordResetTokenExpiry  time.Duration
	SessionExpiry             time.Duration
}

// PlatformConfig holds cross-platform integration settings
type PlatformConfig struct {
	EntativaAPIURL          string
	VignetteAPIURL          string
	EnableCrossPlatformSSO  bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("PORT", "8002"),
			Environment: getEnv("ENV", "development"),
			ServiceName: getEnv("SERVICE_NAME", "vignette-user-service"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "vignette_users"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "vignette-super-secret-key-change-this-in-production"),
			AccessTokenExpiry:   getEnvAsDuration("JWT_ACCESS_TOKEN_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry:  getEnvAsDuration("JWT_REFRESH_TOKEN_EXPIRY", 720*time.Hour),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@vignette.app"),
			FromName:     getEnv("FROM_NAME", "Vignette"),
		},
		Security: SecurityConfig{
			BcryptCost:               getEnvAsInt("BCRYPT_COST", 12),
			PasswordMinLength:        getEnvAsInt("PASSWORD_MIN_LENGTH", 8),
			PasswordResetTokenExpiry: getEnvAsDuration("PASSWORD_RESET_TOKEN_EXPIRY", 1*time.Hour),
			SessionExpiry:            getEnvAsDuration("SESSION_EXPIRY", 24*time.Hour),
		},
		Platform: PlatformConfig{
			EntativaAPIURL:         getEnv("ENTATIVA_API_URL", "http://localhost:8001/api/v1"),
			VignetteAPIURL:         getEnv("VIGNETTE_API_URL", "http://localhost:8002/api/v1"),
			EnableCrossPlatformSSO: getEnvAsBool("ENABLE_CROSS_PLATFORM_SSO", true),
		},
	}
	
	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWT.Secret == "" || c.JWT.Secret == "vignette-super-secret-key-change-this-in-production" {
		if c.Server.Environment == "production" {
			return fmt.Errorf("JWT_SECRET must be set in production")
		}
	}
	return nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}
