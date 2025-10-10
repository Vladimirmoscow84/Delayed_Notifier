package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	RabbitUri    string
	Exchange     string
	Exchangetype string
	Queue        string
	RoutingKey   string
	DLX          string
	DLQ          string
}

type Client struct {
	conn            *amqp.Connection
	publicChannel   *amqp.Channel
	consumerChannel *amqp.Channel
	config          Config
}

func New(cfg Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.RabbitUri)
	if err != nil {
		return nil, fmt.Errorf("[RabbitMQ connection] failed connection to rabbitMq: %w", err)
	}

	publicCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("[RabbitMQ connection] failed to open publicChannel: %w", err)
	}

	consumerCh, err := conn.Channel()
	if err != nil {
		publicCh.Close()
		conn.Close()
		return nil, fmt.Errorf("[RabbitMQ connection] failed to open consumerChannel: %w", err)
	}

	client := &Client{
		conn:            conn,
		publicChannel:   publicCh,
		consumerChannel: consumerCh,
		config:          cfg,
	}
	err = client.setup()
	if err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func (c *Client) setup() error {
	//DLX
	err := c.publicChannel.ExchangeDeclare(
		c.config.DLX,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to declare DLX: %w", err)
	}

	//DLQ
	_, err = c.publicChannel.QueueDeclare(
		c.config.DLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to declare DLQ: %w", err)
	}

	//Bind DLQ to DLX
	err = c.publicChannel.QueueBind(
		c.config.DLQ,
		"",
		c.config.DLX,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to bind DLQ to DLX: %w", err)
	}

	//main exchange
	err = c.publicChannel.ExchangeDeclare(
		c.config.Exchange,
		c.config.Exchangetype,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to declare main exchange: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": c.config.DLX,
	}

	//main queue with DLX
	_, err = c.publicChannel.QueueDeclare(
		c.config.Queue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to declare main queue with DLX: %w", err)
	}

	err = c.publicChannel.QueueBind(
		c.config.Queue,
		c.config.RoutingKey,
		c.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to bind main queue: %w", err)
	}
	return nil
}

type Message struct {
	Body []byte
	TTL  time.Duration
}

func (c *Client) Publish(msg Message) error {
	expiration := fmt.Sprintf("%d", msg.TTL.Milliseconds())
	return c.publicChannel.Publish(
		c.config.Exchange,
		c.config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			Expiration:  expiration,
		},
	)
}

// PublishStruct автоматически сериализует структуру в JSON и публикует сообщение
func (c *Client) PublishStruct(data interface{}, ttl time.Duration) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}
	return c.Publish(Message{Body: body, TTL: ttl})
}

func (c *Client) Ack(msg amqp.Delivery) error {
	return msg.Ack(false)
}

func (c *Client) Nack(msg amqp.Delivery) error {
	return msg.Nack(false, false)
}

func (c *Client) Close() error {
	if c.publicChannel != nil {
		c.publicChannel.Close()
	}
	if c.consumerChannel != nil {
		c.consumerChannel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) ConsumeDLQWithWorkers(ctx context.Context, workerCount int, handler func(msg amqp.Delivery)) error {
	if workerCount <= 0 {
		workerCount = 1
	}

	msgs, err := c.consumerChannel.Consume(
		c.config.DLQ,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume from DLQ: %w", err)
	}

	queue := make(chan amqp.Delivery, workerCount*2)

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case msg, ok := <-queue:
					if !ok {
						return
					}
					handler(msg)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// основной цикл распределения сообщений
	go func() {
		for msg := range msgs {
			select {
			case queue <- msg:
			case <-ctx.Done():
				close(queue)
				return
			}
		}
		close(queue)
	}()

	return nil
}
