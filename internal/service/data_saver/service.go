package datasaver

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
)

type DataSaver struct {
	store store
}

type store interface {
	AddNotice(ctx context.Context, notice model.Notice) (int, error)
}

// type cache interface {
// 	Set(ctx context.Context, key string, value any) error
// }

func New(s store) *DataSaver {
	return &DataSaver{
		store: s,
	}
}
func (s *DataSaver) SaveData(ctx context.Context, notice model.Notice) (int, error) {
	id, err := s.store.AddNotice(ctx, notice)
	if err != nil {
		return 0, fmt.Errorf("[data_saver] failed to add notice: %w", err)
	}
	idStr := strconv.Itoa(id)
	err = s.cache.Set(ctx, idStr, notice.SendStatus)
	if err != nil {
		log.Printf("[data_saver] failed to add cache: %w", err)
	}
	return id, nil
}
