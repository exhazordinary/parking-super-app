package ports

import (
	"context"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

func String(key, value string) Field { return Field{Key: key, Value: value} }
func Err(err error) Field            { return Field{Key: "error", Value: err} }
func Any(key string, val interface{}) Field { return Field{Key: key, Value: val} }

// EventPublisher publishes domain events
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

type Event struct {
	Type    string
	Payload map[string]interface{}
}

const (
	EventProviderCreated     = "provider.created"
	EventProviderActivated   = "provider.activated"
	EventProviderDeactivated = "provider.deactivated"
	EventLocationAdded       = "provider.location.added"
)

// WebhookSender sends webhooks to provider endpoints
type WebhookSender interface {
	Send(ctx context.Context, url string, payload interface{}, secret string) error
}
