package storage

import (
	"context"
	"fmt"
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

// AddNotice - метод, добавляющий уведомление в БД
func (s *Storage) AddNotice(ctx context.Context, notice model.Notice) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO notifications
			(body, date_created, send_date, send_attempts, send_status)
		VALUES
		($1,$2,$3,$4,$5);
		`, notice.Body, notice.DateCreated, notice.SendDate, notice.SendAttempts, notice.SendStatus)
	if err != nil {
		log.Printf("error adding to base %v", err)
		return fmt.Errorf("error adding to base %v", err)
	}
	return nil
}

// DeleteNotice - метод, удаляющий запись из БД
func (s *Storage) DeleteNotice(ctx context.Context, id int) error {
	result, err := s.DB.ExecContext(ctx, `
		DELETE FROM notifications WHERE id=$1;	
	`, id)
	if err != nil {
		log.Printf("error deleting from base %v", err)
		return fmt.Errorf("error deleting from base %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("is not notices with this id %v", err)
		return fmt.Errorf("is not notices with this id %v", err)
	}

	return nil
}
