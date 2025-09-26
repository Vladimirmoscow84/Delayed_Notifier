package model

import "time"

type Notice struct {
	Id           int       `json:"id"`
	Body         string    `json:"body"`
	DateCreated  time.Time `json:"date_created"`
	SendDate     time.Time `json:"send_date"`
	SendAttempts int       `json:"send_attempts"`
	SendStatus   bool      `json:"send_status"`
}
