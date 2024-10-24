package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient implements the RedisAdapter interface
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis adapter
func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	if password != "" {
		rdb.Options().Password = password
	}

	return &RedisClient{client: rdb}
}

// Set stores a key and value in Redis
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	return r.client.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

// Get retrieves a value by key from Redis
func (r *RedisClient) Get(ctx context.Context, key string) (interface{}, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
