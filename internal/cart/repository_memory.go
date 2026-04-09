package cart

type MemoryRepository struct {
	carts map[string]*Cart
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		carts: make(map[string]*Cart),
	}
}

func (r *MemoryRepository) GetByUserID(userID string) (*Cart, error) {
	if cart, exists := r.carts[userID]; exists {
		return cart, nil
	}
	return nil, nil
}

func (r *MemoryRepository) Save(cart *Cart) error {
	r.carts[cart.UserID] = cart
	return nil
}
