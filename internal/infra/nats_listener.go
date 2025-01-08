package infra

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"nats/domain"
	"nats/internal/services"
)

type NATSListener struct {
	js           nats.JetStreamContext
	orderService services.OrderService
}

// NewNATSListener creates a new instance of NATSListener.
func NewNATSListener(js nats.JetStreamContext, orderService services.OrderService) *NATSListener {
	return &NATSListener{
		js:           js,
		orderService: orderService,
	}
}

// StartListening subscribes to the "orders.created" subject.
// Whenever a message arrives, it creates an order using the domain service.
func (l *NATSListener) StartListening() error {
	sub, err := l.js.Subscribe("orders.created", func(msg *nats.Msg) {
		var ord domain.Order
		if err := json.Unmarshal(msg.Data, &ord); err != nil {
			log.Printf("[NATSListener] Failed to unmarshal order: %v\n", err)
			err := msg.Nak()
			if err != nil {
				log.Printf("[NATSListener] Failed to nak order: %v\n", err)
				return
			} // Negative ACK
			return
		}

		err := l.orderService.CreateOrder(context.Background(), ord)
		if err != nil {
			log.Printf("[NATSListener] Failed to create order: %v\n", err)
			err := msg.Nak()
			if err != nil {
				log.Printf("[NATSListener] Failed to nak order: %v\n", err)
				return
			} // or msg.Term() to remove from future deliveries
			return
		}

		log.Printf("[NATSListener] New order created: %+v\n", ord)
		err = msg.Ack()
		if err != nil {
			log.Printf("[NATSListener] Failed to ack order: %v\n", err)
			return
		}
	},
		nats.ManualAck(),
	)
	if err != nil {
		return err
	}

	log.Printf("[NATSListener] Subscribed to subject: %s\n", sub.Subject)
	return nil
}
