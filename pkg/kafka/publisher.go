package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Event represents a domain event to be published
type Event struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
}

// PublisherConfig holds configuration for the Kafka publisher
type PublisherConfig struct {
	Brokers      []string
	Topic        string
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks kafka.RequiredAcks
}

// DefaultPublisherConfig returns sensible default configuration
func DefaultPublisherConfig(brokers []string, topic string) PublisherConfig {
	return PublisherConfig{
		Brokers:      brokers,
		Topic:        topic,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}
}

// Publisher publishes events to Kafka
type Publisher struct {
	writer *kafka.Writer
	tracer trace.Tracer
}

// NewPublisher creates a new Kafka publisher
func NewPublisher(cfg PublisherConfig) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    cfg.BatchSize,
			BatchTimeout: cfg.BatchTimeout,
			RequiredAcks: cfg.RequiredAcks,
		},
		tracer: otel.Tracer("kafka-publisher"),
	}
}

// Publish sends an event to Kafka
func (p *Publisher) Publish(ctx context.Context, event Event) error {
	ctx, span := p.tracer.Start(ctx, "kafka.publish."+event.Type)
	defer span.End()

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Extract trace context
	if spanCtx := trace.SpanFromContext(ctx).SpanContext(); spanCtx.IsValid() {
		event.TraceID = spanCtx.TraceID().String()
		event.SpanID = spanCtx.SpanID().String()
	}

	data, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.Type),
		Value: data,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "timestamp", Value: []byte(event.Timestamp.Format(time.RFC3339))},
		},
	}

	if event.TraceID != "" {
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   "trace_id",
			Value: []byte(event.TraceID),
		})
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// PublishBatch sends multiple events to Kafka
func (p *Publisher) PublishBatch(ctx context.Context, events []Event) error {
	ctx, span := p.tracer.Start(ctx, "kafka.publish.batch")
	defer span.End()

	messages := make([]kafka.Message, len(events))
	for i, event := range events {
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now().UTC()
		}

		if spanCtx := trace.SpanFromContext(ctx).SpanContext(); spanCtx.IsValid() {
			event.TraceID = spanCtx.TraceID().String()
			event.SpanID = spanCtx.SpanID().String()
		}

		data, err := json.Marshal(event)
		if err != nil {
			span.RecordError(err)
			return err
		}

		messages[i] = kafka.Message{
			Key:   []byte(event.Type),
			Value: data,
			Headers: []kafka.Header{
				{Key: "event_type", Value: []byte(event.Type)},
				{Key: "timestamp", Value: []byte(event.Timestamp.Format(time.RFC3339))},
			},
		}
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// Close closes the Kafka writer
func (p *Publisher) Close() error {
	return p.writer.Close()
}

// EventPublisher is the interface that wraps the Publish method
// This matches the ports.EventPublisher interface in services
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

// Ensure Publisher implements EventPublisher
var _ EventPublisher = (*Publisher)(nil)
