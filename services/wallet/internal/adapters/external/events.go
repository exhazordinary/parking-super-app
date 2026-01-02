package external

import (
	"context"
	"encoding/json"
	"log"

	"github.com/parking-super-app/services/wallet/internal/ports"
)

// NoopEventPublisher is a no-op implementation for development.
// In production, replace with Kafka or NATS publisher.
type NoopEventPublisher struct {
	logger *log.Logger
}

func NewNoopEventPublisher() *NoopEventPublisher {
	return &NoopEventPublisher{
		logger: log.Default(),
	}
}

func (p *NoopEventPublisher) Publish(ctx context.Context, event ports.Event) error {
	data, _ := json.Marshal(event.Payload)
	p.logger.Printf("[EVENT] type=%s payload=%s", event.Type, string(data))
	return nil
}

// KafkaEventPublisher publishes events to Kafka.
// This is a placeholder - actual implementation would use kafka-go or sarama.
type KafkaEventPublisher struct {
	brokers []string
	topic   string
}

func NewKafkaEventPublisher(brokers []string, topic string) *KafkaEventPublisher {
	return &KafkaEventPublisher{
		brokers: brokers,
		topic:   topic,
	}
}

func (p *KafkaEventPublisher) Publish(ctx context.Context, event ports.Event) error {
	// Placeholder: In production, use kafka-go writer
	// writer := kafka.NewWriter(kafka.WriterConfig{
	//     Brokers: p.brokers,
	//     Topic:   p.topic,
	// })
	// return writer.WriteMessages(ctx, kafka.Message{
	//     Key:   []byte(event.Type),
	//     Value: payload,
	// })
	log.Printf("[KAFKA] Would publish to topic %s: %s", p.topic, event.Type)
	return nil
}
