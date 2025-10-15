package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8001"),
			Environment:     getEnv("ENVIRONMENT", "development"),
			AllowedOrigins:  []string{getEnv("ALLOWED_ORIGINS", "*")},
			ShutdownTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "socialink"),
			Password:        getEnv("DB_PASSWORD", "socialink_password"),
			DBName:          getEnv("DB_NAME", "socialink_users"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET", "socialink-super-secret-key-change-in-production"),
			AccessTokenTTL:  time.Duration(getEnvAsInt("JWT_ACCESS_TOKEN_TTL", 86400)) * time.Second,  // 24 hours
			RefreshTokenTTL: time.Duration(getEnvAsInt("JWT_REFRESH_TOKEN_TTL", 604800)) * time.Second, // 7 days
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
