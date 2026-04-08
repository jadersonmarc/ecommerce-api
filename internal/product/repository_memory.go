package product

import "errors"

type MemoryRepository struct {
	products map[string]*Product
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
