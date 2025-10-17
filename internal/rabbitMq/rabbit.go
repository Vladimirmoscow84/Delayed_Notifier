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

// PublishStruct автоматически сериализует структуру в JSON и публикует сообщение.
func (c *Client) PublishStructWithTTL(data interface{}, ttl time.Duration) error {
	if err := c.ensurePubChannel(); err != nil {
		return fmt.Errorf("[RabbitMQ]failed to ensure publish channel: %w", err)
	}

	if ttl <= 0 {
		return fmt.Errorf("[RabbitMQ]invalid TTL: must be positive, got %v", ttl)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to marshal struct: %w", err)
	}

	queueName := fmt.Sprintf("msg_%d", time.Now().UnixNano())

	tempQueue, err := c.DeclareTempQueue(queueName, ttl)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to declare temp queue: %w", err)
	}

	if err := c.publicChannel.QueueBind(
		tempQueue,
		tempQueue,
		c.config.Exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("[RabbitMQ]failed to bind temp queue %q: %w", tempQueue, err)
	}

	expiration := fmt.Sprintf("%d", ttl.Milliseconds())
	err = c.publicChannel.Publish(
		c.config.Exchange,
		tempQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  expiration,
		},
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to publish message to temp queue %q: %w", tempQueue, err)
	}

	return nil
}

func (c *Client) ensurePubChannel() error {
	if c.publicChannel != nil && !c.publicChannel.IsClosed() {
		return nil
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to reopen publish channel: %w", err)
	}

	c.publicChannel = ch
	return nil
}

// DeclareTempQueue создает временную очередь с определенным TTL и флагом автоматического удаления.
func (c *Client) DeclareTempQueue(name string, ttl time.Duration) (string, error) {
	ms := ttl.Milliseconds()
	if ms <= 0 {
		return "", fmt.Errorf("[RabbitMQ]invalid TTL: must be positive, got %d ms", ms)
	}

	exp := ms + 5000
	if exp <= 0 {
		return "", fmt.Errorf("[RabbitMQ]invalid x-expires: overflow or negative value (%d)", exp)
	}

	args := amqp.Table{
		"x-message-ttl":          int64(ms),
		"x-expires":              int64(exp),
		"x-dead-letter-exchange": c.config.DLX,
	}

	q, err := c.publicChannel.QueueDeclare(
		name,
		false,
		true,
		false,
		false,
		args,
	)
	if err != nil {
		return "", fmt.Errorf("[RabbitMQ]failed to declare temp queue %q: %w", name, err)
	}

	return q.Name, nil
}

// Ack подтверждения сообщения.
func (c *Client) Ack(msg amqp.Delivery) error {
	return msg.Ack(false)
}

// Nack отрицательное подтверждение сообщения.
func (c *Client) Nack(msg amqp.Delivery) error {
	return msg.Nack(false, false)
}

// Close закрывает оба канала и подключение.
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

// ConsumeDLQWithWorkers использует сообщения из DLQ с рабочим пулом и контекстом для плавного завершения работы.
func (c *Client) ConsumeDLQWithWorkers(ctx context.Context, workerCount int, handler func(msg amqp.Delivery)) error {
	if workerCount <= 0 {
		workerCount = 1
	}

	const consumerTag = "dlq-consumer"

	msgs, err := c.consumerChannel.Consume(
		c.config.DLQ,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("[RabbitMQ]failed to consume from DLQ: %w", err)
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

	go func() {
		<-ctx.Done()
		_ = c.CancelConsumer(consumerTag) // корректное отключение пользователя
		_ = c.consumerChannel.Close()
	}()

	return nil
}

// CancelConsumer корректно отменяет использование пользователя по тегу.
func (c *Client) CancelConsumer(consumerTag string) error {
	if c.consumerChannel == nil {
		return fmt.Errorf("[RabbitMQ]consumer channel is nil")
	}
	if err := c.consumerChannel.Cancel(consumerTag, false); err != nil {
		return fmt.Errorf("[RabbitMQ]failed to cancel consumer: %w", err)
	}
	return nil
}
