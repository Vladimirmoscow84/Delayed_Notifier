package statusgetter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/wb-go/wbf/retry"
)

type cache interface {
	Get(ctx context.Context, strategy retry.Strategy, key string) (string, error)
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
	value, err := s.cache.Get(ctx, str, key)
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
