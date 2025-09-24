package model

import "time"

type Notice struct {
	NoticeUID   string    `json:"notice_uid"`
	Body        string    `json:"body"`
	DateCreated time.Time `json:"date_created"`
	SendDate    time.Time `json:"send_date"`
	SendAtempts int       `json:"send_atempts"`
	SendStatus  bool      `json:"send_status"`
}
