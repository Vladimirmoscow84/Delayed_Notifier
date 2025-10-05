package cache

import (
	"context"
	"time"

	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

type Cache struct {
	client   *redis.Client
	strategy retry.Strategy
}

// New - создание нового кэш
func New(addr string) *Cache {
	client := redis.New(addr, "", 1)
	return &Cache{
		client:   client,
		strategy: retry.Strategy{Attempts: 5, Delay: 3 * time.Second, Backoff: 3},
	}
}

// Get - получение значения из кэш
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.GetWithRetry(ctx, c.strategy, key)
}

// Set - установление значения в кэш
func (c *Cache) Set(ctx context.Context, key string, value any) error {
	return c.client.SetWithRetry(ctx, c.strategy, key, value)
}
