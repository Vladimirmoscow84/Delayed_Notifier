package datasaver

import (
	"context"
	"log"
	"strconv"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	"github.com/wb-go/wbf/retry"
)

type DataSaver struct {
	store store
}

type store interface {
	Set(ctx context.Context, strategy retry.Strategy, key string, value any) error
}

// type cache interface {
// 	Set(ctx context.Context, key string, value any) error
// }

func New(s store) *DataSaver {
	return &DataSaver{
		store: s,
	}
}
func (s *DataSaver) SaveData(ctx context.Context, notice model.Notice) error {
	str := retry.Strategy{}
	idStr := strconv.Itoa(notice.Id)
	err := s.store.Set(ctx, str, idStr, notice.SendStatus)
	if err != nil {
		log.Printf("[data_saver] failed to add cache: %v", err)
	}
	return nil
}
