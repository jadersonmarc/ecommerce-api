package cart

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
)

type Service struct {
	repo           Repository
	productService *product.Service
}

func NewService(r Repository, ps *product.Service) *Service {
	return &Service{repo: r, productService: ps}
}

func generateID() string {
	return uuid.New().String()
}

func (s *Service) AddItems(userID, productID string, quantity int) (*Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}

	product, err := s.productService.GetByID(productID)
	if err != nil {
		return nil, err
	}

	if product.Stock < quantity {
		return nil, errors.New("insufficient stock")
	}

	cart, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = &Cart{
			ID:        generateID(),
			UserID:    userID,
			Items:     []CartItem{},
			CreatedAt: time.Now(),
		}
	}

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     product.Price,
		})
	}

	cart.UpdatedAt = time.Now()

	err = s.repo.Save(cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *Service) GetCart(userID string) (*Cart, error) {
	cart, _ := s.repo.GetByUserID(userID)
	if cart == nil {
		return &Cart{
			UserID: userID,
			Items:  []CartItem{},
		}, nil
	}
	return cart, nil
}

func (s *Service) RemoveItem(userID, productID string) (*Cart, error) {
	cart, _ := s.repo.GetByUserID(userID)
	if cart == nil {
		return nil, errors.New("cart not found")
	}

	var updatedItems []CartItem

	for _, item := range cart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}

	cart.Items = updatedItems
	cart.UpdatedAt = time.Now()

	s.repo.Save(cart)

	return cart, nil
}

func (s *Service) UpdateItem(userID, productId string, quantity int) (*Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}

	cart, _ := s.repo.GetByUserID(userID)
	if cart == nil {
		return nil, errors.New("cart not found")
	}

	for i, item := range cart.Items {
		if item.ProductID == productId {
			product, _ := s.productService.GetByID(productId)

			if product.Stock < quantity {
				return nil, errors.New("insufficient stock")
			}

			cart.Items[i].Quantity = quantity
			cart.UpdatedAt = time.Now()

			s.repo.Save(cart)
			return cart, nil
		}
	}

	return nil, errors.New("item not found in cart")
}

func (s *Service) ClearCart(userID string) error {
	cart, err := s.repo.GetByUserID(userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return nil
	}

	cart.Items = []CartItem{}
	return s.repo.Save(cart)
}
