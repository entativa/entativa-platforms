package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	Security SecurityConfig
	ML       MLConfig
}

type ServerConfig struct {
	Host string
	Port int
	Environment string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type SecurityConfig struct {
	JWTSecret     string
	JWTExpiration int
	BCryptCost    int
}

type MLConfig struct {
	FraudThreshold    float64
	AnomalyThreshold  float64
	TrustScoreVersion string
}

func Load() (*Config, error) {
	// Load .env file if exists
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Port:        getEnvAsInt("SERVER_PORT", 8080),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "meta_users"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", 100),
			MinConns: getEnvAsInt("DB_MIN_CONNS", 10),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 100),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "meta-user-events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "meta-user-service"),
		},
		Security: SecurityConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 3600),
			BCryptCost:    getEnvAsInt("BCRYPT_COST", 10),
		},
		ML: MLConfig{
			FraudThreshold:    getEnvAsFloat("ML_FRAUD_THRESHOLD", 0.7),
			AnomalyThreshold:  getEnvAsFloat("ML_ANOMALY_THRESHOLD", 0.75),
			TrustScoreVersion: getEnv("ML_TRUST_SCORE_VERSION", "v1.0"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
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

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
