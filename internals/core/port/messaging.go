package port

import "time"

// Message defines the message structure
type Message struct {
	Value     []byte
	Key       []byte
	Headers   []MessageHeader
	Timestamp time.Time
	AckFunc   func() error
	NackFunc  func(requeue bool) error
}

// MessageHeader defines the message header structure
type MessageHeader struct {
	Key   string
	Value []byte
}

// MessageHandler is a callback function for processing messages asynchronously
type MessageHandler func(msg Message) error

// Messaging defines the interface for a messaging queue
type Messaging interface {
	Consume(handler MessageHandler) error
	Write(msg Message) error
	Close() error
}
