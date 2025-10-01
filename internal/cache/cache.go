package cache

import (
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
