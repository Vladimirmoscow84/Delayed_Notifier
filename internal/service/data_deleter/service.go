package datadeleter

import (
	"context"
	"fmt"
)

type cache interface {
	Del(ctx context.Context, key string) error
}

type DataDeleter struct {
	cache cache
}

func New(c cache) *DataDeleter {
	return &DataDeleter{cache: c}
}

func (s *DataDeleter) DeleteData(ctx context.Context, key string) error {
	err := s.cache.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("[data_deleter] failed to delete notice from cache by ID: %w", err)
	}
	return nil
}
