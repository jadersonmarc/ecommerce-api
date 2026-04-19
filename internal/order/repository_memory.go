package order

import "errors"

type MemoryRepository struct {
	orders map[string]*Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: make(map[string]*Order),
	}
}

func (r *MemoryRepository) Create(order *Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *MemoryRepository) FindByID(id string) (*Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *MemoryRepository) UpdateStatus(id, status string) error {
	order, exists := r.orders[id]
	if !exists {
		return errors.New("order not found")
	}

	order.Status = status
	return nil
}
