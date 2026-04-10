package product

import (
	"errors"
	"sync"
)

type MemoryRepository struct {
	products map[string]*Product
	mu       sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{products: make(map[string]*Product)}
}

func (r *MemoryRepository) Create(p *Product) error {
	r.products[p.ID] = p
	return nil
}

func (r *MemoryRepository) FindAll() ([]*Product, error) {
	var list []*Product
	for _, p := range r.products {
		list = append(list, p)
	}
	return list, nil
}

func (r *MemoryRepository) FindByID(id string) (*Product, error) {
	p, ok := r.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (r *MemoryRepository) DecreaseStock(productID string, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[productID]
	if !ok {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	product.Stock -= quantity
	return nil
}
