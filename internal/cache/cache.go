package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
)

type Cache struct {
	client   *redis.Client
	strategy retry.Strategy
}

// NewCache - создание нового кэш
func NewCache(client *redis.Client) *Cache {
	fmt.Println("NewCache")
	return &Cache{
		client:   client,
		strategy: retry.Strategy{Attempts: 5, Delay: 3 * time.Second, Backoff: 3},
	}
}

// Get - получение значения из кэш
func (c *Cache) Get(ctx context.Context, key string) (*model.Notice, error) {
	data, err := c.client.GetWithRetry(ctx, c.strategy, key)
	if err != nil {
		return nil, err
	}
	var notice model.Notice
	err = json.Unmarshal([]byte(data), &notice)
	if err != nil {
		return nil, err
	}
	return &notice, nil
}

// Set - установление значения в кэш
func (c *Cache) Set(ctx context.Context, key string, notice *model.Notice) error {
	value, err := json.Marshal(notice)
	if err != nil {
		fmt.Println("error marshaling data:", err)
		return err
	}
	return c.client.SetWithRetry(ctx, c.strategy, key, value)
}
