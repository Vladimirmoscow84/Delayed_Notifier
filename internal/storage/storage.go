package storage

import (
	"context"
	"log"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// структура для работы с БД
type Storage struct {
	DB *sqlx.DB
}

// New - конструктор для создания экземпляра Storage
func New(databaseUri string) (*Storage, error) {
	db, err := sqlx.Connect("pgx", databaseUri)
	if err != nil {
		log.Fatalf("connection to DB error %v", err)
	}
	return &Storage{
		DB: db,
	}, nil
}

// AddNotification - метод, добавляющий уведомление в БД
func (s *Storage) AddNotice(ctx context.Context, notice model.Notice) error {

	return nil
}
