package order

type Repository interface {
	Create(order *Order) error
	FindByID(id string) (*Order, error)
}
