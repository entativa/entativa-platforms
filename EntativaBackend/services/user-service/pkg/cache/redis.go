package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"socialink/user-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with expiration
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(key string, dest interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key from Redis
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisCache) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

// Increment increments a counter
func (r *RedisCache) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// IncrementWithExpiry increments a counter and sets expiry if new
func (r *RedisCache) IncrementWithExpiry(key string, expiration time.Duration) (int64, error) {
	pipe := r.client.Pipeline()
	incr := pipe.Incr(r.ctx, key)
	pipe.Expire(r.ctx, key, expiration)
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// SetNX sets a key only if it doesn't exist (returns true if set)
func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.SetNX(r.ctx, key, data, expiration).Result()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}
