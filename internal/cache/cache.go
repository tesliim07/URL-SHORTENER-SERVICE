package cache

import (
	"context"
	"fmt"
	"url-shortener-service/config"
	"github.com/redis/go-redis/v9"
)


type URLCache interface {
	SetURL(code, originalURL string) error
	GetURL(code string) (string, error)
}

// Cache is the Redis implementation of URLCache
type Cache struct {
	client *redis.Client
}

// NewCache creates a new Cache and connects to Redis
func NewCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	// ping redis to check connection is alive
	ctx := context.Background() //requires a context on every operation and context.Background() just means no timeout and no cancellation
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// SetURL stores a short code → original URL in Redis
func (c *Cache) SetURL(code, originalURL string) error {
	ctx := context.Background()
	err := c.client.Set(ctx, code, originalURL, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to cache url: %w", err)
	}
	return nil
}

// GetURL retrieves a cached original URL by its short code
func (c *Cache) GetURL(code string) (string, error) {
	ctx := context.Background()
	originalURL, err := c.client.Get(ctx, code).Result()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve cached url: %w", err)
	}
	return originalURL, nil
}