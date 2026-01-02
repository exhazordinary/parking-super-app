package external

import (
	"context"
	"encoding/json"
	"log"

	"github.com/parking-super-app/services/provider/internal/ports"
)

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
