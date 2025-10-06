package statusgetter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
)

type StatusGetter struct {
	cache cache
}

type cache interface {
	Get(ctx context.Context, key string) (string, error)
}

func New(c cache) *StatusGetter {
	return &StatusGetter{cache: c}
}

func (s *StatusGetter) GetStatusNotice(ctx context.Context, key string) (string, error) {
	value, err := s.cache.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to get notice from cache: %w", err)
	}
	var notice model.Notice

	err = json.Unmarshal([]byte(value), &notice)
	if err != nil {
		return "", fmt.Errorf("[status_getter] failed to unmarshal notice: %w", err)
	}

	return notice.SendStatus, nil
}
