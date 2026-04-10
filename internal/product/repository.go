package product

type Repository interface {
	Create(product *Product) error
	FindAll() ([]*Product, error)
	FindByID(id string) (*Product, error)
	DecreaseStock(id string, quantity int) error
}
