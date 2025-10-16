package email

import (
	"bytes"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/Vladimirmoscow84/Delayed_Notifier.git/internal/model"
)

type Client struct {
	host string
	port int
	user string
	pass string
	from string
	to   []string
}

func New(host, portStr, user, pass, from, toList string) *Client {
	if host == "" || portStr == "" || user == "" || pass == "" || from == "" || toList == "" {
		log.Println("[email]: missing EMAIL config")
		return nil
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("[email] invalid EMAIL_PORT")
		return nil
	}

	return &Client{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
		to:   strings.Split(toList, ","),
	}
}

// Send - отправка уведомления по email
func (c *Client) Send(notice model.Notice) {
	if c == nil {
		return
	}
	subject := "Notification ID=" + strconv.Itoa(notice.Id)
	body := "Уведомление ID=" + strconv.Itoa(notice.Id) + "\n\n" +
		notice.Body + "\n\nЗапланировано: " + notice.SendDate.Format(time.RFC3339)

	msg := bytes.Buffer{}
	msg.WriteString("From: " + c.from + "\r\n")
	msg.WriteString("To: " + strings.Join(c.to, ",") + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", c.user, c.pass, c.host)
	addr := c.host + ":" + strconv.Itoa(c.port)

	maxAttempts := notice.SendAttempts
	backoff := time.Second * 1
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := smtp.SendMail(addr, auth, c.from, c.to, msg.Bytes())
		if err == nil {
			log.Printf("[Email] notice ID=%d sent", notice.Id)
			return
		}
		log.Printf("[Email] attempt %d failed: %v", attempt, err)
		time.Sleep(backoff)
		backoff *= 2
	}
	log.Printf("[Email] notice ID=%d failed after %d attempts", notice.Id, maxAttempts)
}
