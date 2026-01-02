package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// EventHandler is a function that handles a specific event type
type EventHandler func(ctx context.Context, event Event) error

// ConsumerConfig holds configuration for the Kafka consumer
type ConsumerConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	MinBytes int
	MaxBytes int
}

// DefaultConsumerConfig returns sensible default configuration
func DefaultConsumerConfig(brokers []string, topic, groupID string) ConsumerConfig {
	return ConsumerConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}
}

// Consumer consumes events from Kafka
type Consumer struct {
	reader   *kafka.Reader
	handlers map[string]EventHandler
	mu       sync.RWMutex
	tracer   trace.Tracer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg ConsumerConfig) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			Topic:    cfg.Topic,
			GroupID:  cfg.GroupID,
			MinBytes: cfg.MinBytes,
			MaxBytes: cfg.MaxBytes,
		}),
		handlers: make(map[string]EventHandler),
		tracer:   otel.Tracer("kafka-consumer"),
	}
}

// RegisterHandler registers a handler for a specific event type
func (c *Consumer) RegisterHandler(eventType string, handler EventHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[eventType] = handler
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return c.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil // Context cancelled
				}
				log.Printf("error fetching message: %v", err)
				continue
			}

			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("error processing message: %v", err)
				// Continue processing even if one message fails
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var event Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("error unmarshaling event: %v", err)
		return err
	}

	ctx, span := c.tracer.Start(ctx, "kafka.consume."+event.Type)
	defer span.End()

	c.mu.RLock()
	handler, ok := c.handlers[event.Type]
	c.mu.RUnlock()

	if !ok {
		// No handler registered for this event type, skip
		return nil
	}

	if err := handler(ctx, event); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// Close closes the Kafka reader
func (c *Consumer) Close() error {
	return c.reader.Close()
}
