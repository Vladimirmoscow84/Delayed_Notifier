package statusgetter

import (
	"context"
	"fmt"
	"log"

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
	}

	exists, err := s.cache.Exists(ctx, key)
	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to check cache key: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("[status_getter] key %q not found in cache", key)
	}

	value, err := s.cache.Get(ctx, str, key)
	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to get notice from cache: %w", err)
	}

	log.Printf("[status_getter] got cached status for key %q: %v", key, value)

	return value, nil
}
