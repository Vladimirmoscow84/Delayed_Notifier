package datadeleter

import (
	"context"
	"fmt"
)

type cache interface {
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type DataDeleter struct {
	cache cache
}

func New(c cache) *DataDeleter {
	return &DataDeleter{cache: c}
}

func (s *DataDeleter) DeleteData(ctx context.Context, key string) error {
	exists, err := s.cache.Exists(ctx, key)
	if err != nil {
		return fmt.Errorf("[data_deleter] failed to check cache key: %w", err)
	}
	if !exists {
		return fmt.Errorf("[status_deleter] key %q not found in cache", key)
	}
	err = s.cache.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("[data_deleter] failed to delete notice from cache by ID: %w", err)
	}
	return nil
}
