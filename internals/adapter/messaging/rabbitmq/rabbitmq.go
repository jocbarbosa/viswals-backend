package rabbitmq

import (
	"fmt"

	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"github.com/streadway/amqp"
)

// RabbitMQAdapter implements the Messaging interface for RabbitMQ
type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

// NewRabbitMQAdapter creates a new RabbitMQ adapter
func NewRabbitMQAdapter(connURL string, queueName string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQAdapter{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

// Write publishes a message to RabbitMQ
func (r *RabbitMQAdapter) Write(msg port.Message) error {
	err := r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg.Value,
			Headers:     r.convertHeaders(msg.Headers),
			Timestamp:   msg.Timestamp,
		})
	if err != nil {
		return err
	}

	return nil
}

// Consume starts consuming messages from the queue and calls the provided handler
func (r *RabbitMQAdapter) Consume(handler port.MessageHandler) error {
	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for d := range msgs {
			var msg port.Message

			msg.Value = d.Body
			msg.Headers = convertAMQPHeaders(d.Headers)
			msg.Timestamp = d.Timestamp

			err := handler(msg)
			if err != nil {
				fmt.Println("error processing message:", err)
			} else {
				err = d.Ack(false)
				if err != nil {
					fmt.Println("error acknowledging message:", err)
				}
			}
		}
	}()

	return nil
}

// Close closes the connection to RabbitMQ
func (r *RabbitMQAdapter) Close() error {
	if err := r.connection.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	return r.connection.Close()
}

func (r *RabbitMQAdapter) convertHeaders(headers []port.MessageHeader) amqp.Table {
	table := amqp.Table{}
	for _, header := range headers {
		table[header.Key] = header.Value
	}
	return table
}

func convertAMQPHeaders(headers amqp.Table) []port.MessageHeader {
	convertedHeaders := make([]port.MessageHeader, 0, len(headers))
	for key, value := range headers {
		strVal, ok := value.(string)
		if !ok {
			continue
		}
		convertedHeaders = append(convertedHeaders, port.MessageHeader{Key: key, Value: []byte(strVal)})
	}
	return convertedHeaders
}
