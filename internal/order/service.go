package order

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
)

type PaymentService interface {
	CreatePaymentIntent(amount int64, orderID string) (string, error)
}
type Service struct {
	repo           Repository
	cartService    *cart.Service
	productService *product.Service
	paymentService PaymentService
}
type CheckoutResponse struct {
	Order        *Order
	ClientSecret string
}

func GenerateID() string {
	return uuid.New().String()
}

func NewService(r Repository, cs *cart.Service, ps *product.Service, pay PaymentService) *Service {
	return &Service{
		repo:           r,
		cartService:    cs,
		productService: ps,
		paymentService: pay,
	}
}

func (s *Service) Checkout(userID string) (*CheckoutResponse, error) {
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

	ClientSecret, err := s.paymentService.CreatePaymentIntent(total, order.ID)
	if err != nil {
		return nil, err
	}

	return &CheckoutResponse{
		Order:        order,
		ClientSecret: ClientSecret,
	}, nil

}

func (s *Service) Pay(orderID string) error {
	order, err := s.repo.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.Status == StatusPaid {
		return errors.New("order already paid")
	}

	for _, item := range order.Items {
		err := s.productService.DecreaseStock(item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	if err := s.repo.UpdateStatus(orderID, StatusPaid); err != nil {
		return err
	}

	err = s.cartService.ClearCart(order.UserID)
	if err != nil {
		return err
	}

	return nil
}
