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

// NewCache - создание нового кэш
func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client:   client,
		strategy: retry.Strategy{Attempts: 5, Delay: 3 * time.Second, Backoff: 3},
	}
}

// Get - получение значения из кэш
func (c *Cache) Get(ctx context.Context, strategy retry.Strategy, key string) (string, error) {
	return c.client.GetWithRetry(ctx, c.strategy, key)
}

// Set - установление значения в кэш
func (c *Cache) Set(ctx context.Context, strategy retry.Strategy, key string, value any) error {
	return c.client.SetWithRetry(ctx, c.strategy, key, value)
}

// Del - удаляет значение из кэш
func (c *Cache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key)
}

// Exists - проверяет наличие записей в кэш по ключу
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
