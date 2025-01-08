package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nats/domain"

	"github.com/nats-io/nats.go"
)

// NATSPublisher is the outbound adapter for NATS JetStream.
type NATSPublisher struct {
	js nats.JetStreamContext
}

// NewNATSPublisher creates a new NATSPublisher.
func NewNATSPublisher(js nats.JetStreamContext) *NATSPublisher {
	return &NATSPublisher{
		js: js,
	}
}

// PublishOrderProcessed publishes an event with the updated order to "orders.processed".
func (p *NATSPublisher) PublishOrderProcessed(ctx context.Context, order domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	ack, err := p.js.Publish("orders.processed", data)
	if err != nil {
		return fmt.Errorf("failed to publish order processed event: %w", err)
	}

	log.Printf("[NATSPublisher] 'orders.processed' event published. Seq=%d Stream=%s\n", ack.Sequence, ack.Stream)
	return nil
}
