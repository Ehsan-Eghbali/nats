package repository

import (
	"context"
	"errors"
	"fmt"
	"nats/domain"
	"sync"
)

// MemoryOrderRepository is an in-memory implementation of OrderRepository.
type MemoryOrderRepository struct {
	data map[string]domain.Order
	mu   sync.RWMutex
}

// NewMemoryOrderRepository creates a new in-memory repository.
func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		data: make(map[string]domain.Order),
	}
}

func (r *MemoryOrderRepository) Save(ctx context.Context, order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[order.ID]; exists {
		return fmt.Errorf("order with ID=%s already exists", order.ID)
	}
	r.data[order.ID] = order
	return nil
}

func (r *MemoryOrderRepository) FindByID(ctx context.Context, id string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ord, ok := r.data[id]
	if !ok {
		return domain.Order{}, errors.New("order not found")
	}
	return ord, nil
}

func (r *MemoryOrderRepository) Update(ctx context.Context, order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[order.ID]; !exists {
		return errors.New("order not found for update")
	}
	r.data[order.ID] = order
	return nil
}
