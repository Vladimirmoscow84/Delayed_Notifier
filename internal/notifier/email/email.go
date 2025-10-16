package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
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
func (c *Client) Send(subject, body string) error {
	if c == nil {
		return fmt.Errorf("[email] client is empty")
	}
	if len(c.to) == 0 {
		return fmt.Errorf("[email] no recipients in EMAIL_TO")
	}

	auth := smtp.PlainAuth("", c.user, c.pass, c.host)
	msg := fmt.Sprintf("From: %s\r\n", c.from)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(c.to, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-version: 1.0;\r\n"
	msg += "Content-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"
	msg += body

	addr := c.host + ":" + strconv.Itoa(c.port)

	backoff := 1 * time.Second
	for attempt := 1; attempt <= 3; attempt++ {
		err := smtp.SendMail(addr, auth, c.from, c.to, []byte(msg))
		if err == nil {
			log.Printf("[Email] successfully sended")
			return nil
		}
		log.Printf("[Email] attempt %d failed: %v", attempt, err)
		time.Sleep(backoff)
		backoff *= 2
	}
	log.Printf("[Email] send failed after all attempts")
	return fmt.Errorf("[email] failed send to email")
}
