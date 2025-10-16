package statusgetter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wb-go/wbf/retry"
)

type cache interface {
	Get(ctx context.Context, strategy retry.Strategy, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}
type StatusGetter struct {
	cache cache
}

func New(c cache) *StatusGetter {
	return &StatusGetter{cache: c}
}

func (s *StatusGetter) GetStatusNotice(ctx context.Context, key string) (string, error) {
	str := retry.Strategy{
		Attempts: 5,
		Delay:    300 * time.Millisecond,
	}
	var exists bool
	var value string
	var err error
	//проверка наличия ключа в соответстви со стратегией попыток
	for i := 0; i < str.Attempts; i++ {
		exists, err = s.cache.Exists(ctx, key)
		if err == nil {
			break
		}
		log.Printf("[status_getter] attempt %d to check key from cache failed: %v", i)
		time.Sleep(str.Delay)
	}

	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to check cache key after %d attempts: %w", str.Attempts, err)
	}
	if !exists {
		return "", fmt.Errorf("[status_getter] key %q not found in cache", key)
	}

	//получение значений из кэш по ключу в соответствии со стратегией
	for i := 0; i < str.Attempts; i++ {
		value, err = s.cache.Get(ctx, str, key)
		if err == nil {
			break
		}
		log.Printf("[status_getter] attempt %d to get key from cache failed: %v", i+1, err)
		time.Sleep(str.Delay)
	}

	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to get notice from cache afte %d attempts: %w", str.Attempts, err)
	}

	log.Printf("[status_getter] got cached status for key %q: %v", key, value)

	return value, nil
}
