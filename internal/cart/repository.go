package cart

type Repository interface {
	GetByUserID(userID string) (*Cart, error)
	Save(cart *Cart) error
}
