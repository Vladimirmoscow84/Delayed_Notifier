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
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows count: %v", err)
		return fmt.Errorf("error retrieving affected row count: %w", err)
	}
	if rowsAffected == 0 {
		log.Printf("is not notices with this id = %v", id)
		return fmt.Errorf("is not notices with this id = %v", id)
	}

	return nil
}

// GetNotice - метод для получения одной записи из БД
func (s *Storage) GetNotice(ctx context.Context, id int) (model.Notice, error) {
	var notice model.Notice

	row := s.DB.QueryRowContext(ctx, `
	SELECT id, body, date_created, send_date, send_attempts, send_status 
	FROM notifications WHERE id = $1;
	`, id)
	err := row.Scan(notice.Id, notice.Body, notice.DateCreated, notice.SendDate, notice.SendAttempts, notice.SendStatus)
	if err != nil {
		log.Printf("get notice error %v", err)
		return model.Notice{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return notice, nil
}

// GetNiticies - метод для получения всех записей из БД{
func (s *Storage) GetNoticies(ctx context.Context) ([]model.Notice, error) {
	var notics []model.Notice

	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, body, date_created, send_date, send_attempts, send_status FROM notifications
	`)
	if err != nil {
		log.Printf("get notices error %v", err)
		return []model.Notice{}, fmt.Errorf("failed to scan rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notice model.Notice
		err := rows.Scan(&notice.Id, &notice.Body, &notice.DateCreated, &notice.SendDate, &notice.SendAttempts, &notice.SendStatus)
		if err != nil {
			log.Printf("scan row error %v", err)
			continue //пропускаем ошибку
		}
		notics = append(notics, notice)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows: %w", err)
	}

	return notics, nil
}
