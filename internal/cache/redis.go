package cache

import (
	"context"
	"encoding/json"
	// "errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache handles all Redis caching operations
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Set stores a value in Redis with TTL
// value can be any struct that is JSON serializable
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return rc.client.Set(ctx, key, string(data), ttl).Err()
}

// Get retrieves a value from Redis and unmarshals it into dest
// Returns redis.Nil error if key doesn't exist (cache miss)
// Returns other errors if there are Redis connection or unmarshal issues
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		// Return the error directly (including redis.Nil for cache miss)
		return err
	}

	// Only unmarshal if we have a value
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// GetString retrieves a string value from Redis
// Returns redis.Nil error if key doesn't exist (cache miss)
func (rc *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		// Return the error directly (including redis.Nil for cache miss)
		return "", err
	}
	return val, nil
}

// SetString stores a string value in Redis with TTL
func (rc *RedisCache) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return rc.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from Redis
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// DeleteMultiple removes multiple keys from Redis
func (rc *RedisCache) DeleteMultiple(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return rc.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists in Redis
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return count > 0, nil
}

// InvalidatePattern deletes all keys matching a pattern
// Example: InvalidatePattern(ctx, "product:*") deletes all product keys
func (rc *RedisCache) InvalidatePattern(ctx context.Context, pattern string) error {
	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		return rc.client.Del(ctx, keys...).Err()
	}

	return nil
}

// SetWithCondition sets a value only if the key doesn't exist (NX = Not eXists)
func (rc *RedisCache) SetWithCondition(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := rc.client.SetNX(ctx, key, string(data), ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set with condition: %w", err)
	}

	return result, nil
}

// Expire sets expiration time for an existing key
func (rc *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return rc.client.Expire(ctx, key, ttl).Err()
}

// TTL returns the remaining time to live of a key (-1 if no expiry, -2 if doesn't exist)
func (rc *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := rc.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}
	return ttl, nil
}

// Increment increments the integer value of a key by 1
func (rc *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	val, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key: %w", err)
	}
	return val, nil
}

// IncrementBy increments the integer value of a key by the given amount
func (rc *RedisCache) IncrementBy(ctx context.Context, key string, increment int64) (int64, error) {
	val, err := rc.client.IncrBy(ctx, key, increment).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key: %w", err)
	}
	return val, nil
}

// Decrement decrements the integer value of a key by 1
func (rc *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := rc.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key: %w", err)
	}
	return val, nil
}

// DecrementBy decrements the integer value of a key by the given amount
func (rc *RedisCache) DecrementBy(ctx context.Context, key string, decrement int64) (int64, error) {
	val, err := rc.client.DecrBy(ctx, key, decrement).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key: %w", err)
	}
	return val, nil
}

// FlushAll deletes all keys in the current Redis database (USE WITH CAUTION!)
func (rc *RedisCache) FlushAll(ctx context.Context) error {
	return rc.client.FlushAll(ctx).Err()
}

// Close closes the Redis client connection
func (rc *RedisCache) Close() error {
	if rc.client != nil {
		return rc.client.Close()
	}
	return nil
}

// Health checks Redis connection health
func (rc *RedisCache) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return rc.client.Ping(ctx).Err()
}