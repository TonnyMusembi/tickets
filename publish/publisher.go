package publish

import (
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var ch *amqp091.Channel

func InitRabbitMQ(url string) error {
	var err error
	conn, err = amqp091.Dial(url)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func Publish(queueName string, body []byte) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		slog.Error("Failed to declare queue", "error", err)
		return err
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	slog.Info("Published message to queue", "queue", queueName, "body", string(body))

	return err
}

func Consume(queueName string) (<-chan amqp091.Delivery, error) {
	_, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		slog.Error("Failed to declare queue", "error", err)
		return nil, err
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		slog.Error("Failed to register consumer", "error", err)
		return nil, err
	}

	return msgs, nil
}
