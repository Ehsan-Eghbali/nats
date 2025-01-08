package main

import (
	"context"
	"log"
	"nats/config"
	"nats/internal/infra"
	"nats/internal/repository"
	"nats/internal/services"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to NATS
	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("[Main] Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 3. Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("[Main] Failed to create JetStream context: %v", err)
	}

	// 4. Ensure the stream exists (or create/update it)
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.StreamName,
		Subjects: []string{"orders.*"},
		Storage:  nats.FileStorage, // or nats.MemoryStorage
	})
	if err != nil {
		log.Printf("[Main] Warning: could not Add/Update stream: %v", err)
	}

	// 5. Initialize the repository
	orderRepo := repository.NewMemoryOrderRepository()

	// 6. Initialize the domain service
	orderService := services.NewOrderService(orderRepo)

	// 7. Initialize the inbound adapter (Listener)
	listener := infra.NewNATSListener(js, orderService)
	if err := listener.StartListening(); err != nil {
		log.Fatalf("[Main] Failed to start NATS listener: %v", err)
	}

	// 8. Initialize the outbound adapter (Publisher)
	publisher := infra.NewNATSPublisher(js)

	// 9. Optional demonstration: after a few seconds, we process an existing order and publish an event
	go func() {
		// Wait for a few seconds so we can first publish an "orders.created" message externally
		time.Sleep(5 * time.Second)

		// In a real scenario, we might dynamically pick an order ID, etc.
		// Let's assume "order-101" was previously published to "orders.created".
		err := orderService.ProcessOrder(context.Background(), "order-101")
		if err != nil {
			log.Printf("[Main] ProcessOrder error: %v", err)
			return
		}

		// Retrieve the updated order
		ord, _ := orderService.GetOrder(context.Background(), "order-101")

		// Publish the "processed" event
		err = publisher.PublishOrderProcessed(context.Background(), ord)
		if err != nil {
			log.Printf("[Main] PublishOrderProcessed error: %v", err)
		}
	}()

	log.Println("[Main] Service is running. Press Ctrl+C to exit.")
	select {}
}
