package order

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
)

type Service struct {
	repo           Repository
	cartService    *cart.Service
	productService *product.Service
}

func GenerateID() string {
	return uuid.New().String()
}

func NewService(r Repository, cs *cart.Service, ps *product.Service) *Service {
	return &Service{
		repo:           r,
		cartService:    cs,
		productService: ps,
	}
}

func (s *Service) Checkout(userID string) (*Order, error) {
	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("empty cart")
	}

	var total int64
	var items []OrderItem

	for _, item := range cart.Items {
		product, err := s.productService.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product " + product.Name)
		}

		items = append(items, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		total += int64(item.Quantity) * product.Price

	}

	order := &Order{
		ID:        GenerateID(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) Pay(orderID string) error {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.Status == StatusPaid {
		return errors.New("order already paid")
	}

	order.Status = StatusPaid

	for _, item := range order.Items {
		product, err := s.productService.GetByID(item.ProductID)
		if err != nil {
			return err
		}
		product.Stock -= item.Quantity
	}
	err = s.cartService.ClearCart(order.UserID)
	if err != nil {
		return err
	}

	return nil
}
