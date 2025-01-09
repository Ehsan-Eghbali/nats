package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

// A simple in-memory store for checking idempotency.
var processedIDs = struct {
	sync.RWMutex
	data map[string]bool
}{
	data: make(map[string]bool),
}

func main() {
	// 1) Connect to NATS (JetStream enabled).
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 2) Create the JetStream context.
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to get JetStream context: %v", err)
	}

	// 3) Create (or update) a main stream for orders.
	//    We store any subject matching "orders.*".
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS_STREAM",
		Subjects: []string{"orders.*"},
		Storage:  nats.FileStorage, // or nats.MemoryStorage
	})
	if err != nil {
		log.Printf("AddStream ORDERS_STREAM error (maybe exists): %v", err)
	}

	// 4) Create (or update) a separate stream for DLQ (optional).
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS_DLQ",
		Subjects: []string{"dlq.orders.*"},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Printf("AddStream ORDERS_DLQ error: %v", err)
	}

	// 5) Subscribe to "orders.created" via JetStream in a queue group.
	//    We specify nats.ManualAck() so we can ack manually.
	sub, err := js.QueueSubscribe(
		"orders.created", // Subject
		"orders-workers", // Queue group
		func(msg *nats.Msg) {
			handleOrderMessage(msg)
		},
		nats.ManualAck(),
	)
	if err != nil {
		log.Fatalf("Failed to queue-subscribe via JetStream: %v", err)
	}

	// 6) As an example, publish a few messages in a separate goroutine.
	go func() {
		for i := 1; i <= 5; i++ {
			orderID := fmt.Sprintf("order-%d", i)
			if _, err := js.Publish("orders.created", []byte(orderID)); err != nil {
				log.Printf("Publish error: %v", err)
			} else {
				log.Printf("Published new order: %s", orderID)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	log.Println("JetStream queue subscription is running... (Press Ctrl+C to exit)")
	select {} // Block forever
	_ = sub
}

// handleOrderMessage applies idempotency checks and simulates an error for "order-3".
func handleOrderMessage(msg *nats.Msg) {
	id := string(msg.Data)

	// Idempotency check.
	if alreadyProcessed(id) {
		log.Printf("[Duplicate] Already processed %s, just ack", id)
		_ = msg.Ack()
		return
	}
	markProcessed(id)

	log.Printf("Processing order: %s", id)

	// Simulate an error for order-3 to test re-delivery.
	if id == "order-3" {
		log.Printf("[Error] Failed to process %s, no Ack -> it will be re-delivered", id)
		return
	}

	// Ack if successful.
	_ = msg.Ack()
	log.Printf("Acked %s successfully", id)
}

// alreadyProcessed returns true if we've seen this message ID before.
func alreadyProcessed(id string) bool {
	processedIDs.RLock()
	defer processedIDs.RUnlock()
	return processedIDs.data[id]
}

// markProcessed records this ID so we don't process it again.
func markProcessed(id string) {
	processedIDs.Lock()
	defer processedIDs.Unlock()
	processedIDs.data[id] = true
}
