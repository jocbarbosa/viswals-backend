package rabbitmq

import (
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"github.com/streadway/amqp"
)

type RabbitMQAdapter struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

// NewRabbitMQAdapter creates a new instance of RabbitMQAdapter
func NewRabbitMQAdapter(url, queueName string) (*RabbitMQAdapter, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQAdapter{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

// Write sends a message to RabbitMQ
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

// Consume sets up a consumer to handle messages asynchronously
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
		return err
	}

	go func() {
		for d := range msgs {
			message := port.Message{
				Value:     d.Body,
				Headers:   r.convertAmqpHeaders(d.Headers),
				Timestamp: d.Timestamp,
				AckFunc: func() error {
					return d.Ack(false)
				},
				NackFunc: func(requeue bool) error {
					return d.Nack(false, requeue)
				},
			}

			err := handler(message)
			if err != nil {
				_ = message.NackFunc(true)
			} else {
				_ = message.AckFunc()
			}
		}
	}()

	return nil
}

// Close closes the connection and channel
func (r *RabbitMQAdapter) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.connection.Close()
}

// convertHeaders converts custom headers to amqp.Table format
func (r *RabbitMQAdapter) convertHeaders(headers []port.MessageHeader) amqp.Table {
	table := amqp.Table{}
	for _, header := range headers {
		table[header.Key] = header.Value
	}
	return table
}

func (r *RabbitMQAdapter) convertAmqpHeaders(headers amqp.Table) []port.MessageHeader {
	var result []port.MessageHeader
	for key, value := range headers {
		if val, ok := value.([]byte); ok {
			result = append(result, port.MessageHeader{
				Key:   key,
				Value: val,
			})
		}
	}
	return result
}
