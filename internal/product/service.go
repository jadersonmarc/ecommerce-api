package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func generateID() string {
	return uuid.New().String()
}

func (s *Service) Create(name, description string, price int64, stock int) (*Product, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if price <= 0 {
		return nil, errors.New("invalid price")
	}

	if stock < 0 {
		return nil, errors.New("invalid stock")
	}

	product := &Product{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		createdAt:   time.Now().Unix(),
	}

	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil

}

func (s *Service) List() ([]*Product, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id string) (*Product, error) {
	return s.repo.FindByID(id)
}
