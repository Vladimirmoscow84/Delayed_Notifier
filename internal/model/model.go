package model

import "time"

type Notice struct {
	Id           int       `json:"id"`
	Body         string    `json:"body" binding:"required"`
	DateCreated  time.Time `json:"date_created"`
	SendDate     time.Time `json:"send_date" binding:"required"`
	SendAttempts int       `json:"send_attempts"`
	SendStatus   string    `json:"send_status"`
}
