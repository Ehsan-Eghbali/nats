package services

import (
	"context"
	"errors"
	"fmt"
	"nats/domain"
)

// OrderService defines the domain service methods for managing orders.
type OrderService interface {
	CreateOrder(ctx context.Context, order domain.Order) error
	GetOrder(ctx context.Context, id string) (domain.Order, error)
	ProcessOrder(ctx context.Context, id string) error
}

// OrderRepository is an interface for accessing order data (in hex-architecture, this is the "driven port").
type OrderRepository interface {
	Save(ctx context.Context, order domain.Order) error
	FindByID(ctx context.Context, id string) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
}

// orderService is our default implementation of OrderService.
type orderService struct {
	repo OrderRepository
}

// NewOrderService creates a new instance of our domain service with the given repository.
func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{
		repo: repo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order domain.Order) error {
	if order.ID == "" {
		return errors.New("order ID must not be empty")
	}
	// Check if order already exists
	existing, err := s.repo.FindByID(ctx, order.ID)
	if err == nil && existing.ID != "" {
		return fmt.Errorf("order with ID=%s already exists", order.ID)
	}

	order.Status = "created"
	return s.repo.Save(ctx, order)
}

func (s *orderService) GetOrder(ctx context.Context, id string) (domain.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *orderService) ProcessOrder(ctx context.Context, id string) error {
	ord, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find order with ID=%s: %w", id, err)
	}
	if ord.Status != "created" {
		return fmt.Errorf("cannot process order with current status: %s", ord.Status)
	}
	ord.Status = "processed"
	return s.repo.Update(ctx, ord)
}
