package order

import (
	"testing"

	"github.com/jadersonmarc/ecommerce-api/internal/cart"
	"github.com/jadersonmarc/ecommerce-api/internal/product"
)

type paymentServiceStub struct {
	clientSecret string
	amount       int64
	orderID      string
	calls        int
	err          error
}

func (s *paymentServiceStub) CreatePaymentIntent(amount int64, orderID string) (string, error) {
	s.calls++
	s.amount = amount
	s.orderID = orderID
	if s.err != nil {
		return "", s.err
	}
	return s.clientSecret, nil
}

func TestServiceCheckoutCreatesPendingOrderAndPaymentIntent(t *testing.T) {
	productRepo := product.NewMemoryRepository()
	productService := product.NewService(productRepo)
	createdProduct, err := productService.Create("Keyboard", "Mechanical keyboard", 2500, 10)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	cartRepo := cart.NewMemoryRepository()
	cartService := cart.NewService(cartRepo, productService)
	if _, err := cartService.AddItems("user-1", createdProduct.ID, 2); err != nil {
		t.Fatalf("AddItems() error = %v", err)
	}

	paymentStub := &paymentServiceStub{clientSecret: "pi_secret_test"}
	orderService := NewService(NewMemoryRepository(), cartService, productService, paymentStub)

	response, err := orderService.Checkout("user-1")
	if err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}

	if response.ClientSecret != "pi_secret_test" {
		t.Fatalf("expected client secret %q, got %q", "pi_secret_test", response.ClientSecret)
	}

	if response.Order == nil {
		t.Fatal("expected created order in response")
	}

	if response.Order.Status != StatusPending {
		t.Fatalf("expected order status %q, got %q", StatusPending, response.Order.Status)
	}

	if response.Order.Total != 5000 {
		t.Fatalf("expected total 5000, got %d", response.Order.Total)
	}

	if len(response.Order.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(response.Order.Items))
	}

	if paymentStub.calls != 1 {
		t.Fatalf("expected payment service to be called once, got %d", paymentStub.calls)
	}

	if paymentStub.amount != 5000 {
		t.Fatalf("expected payment amount 5000, got %d", paymentStub.amount)
	}

	if paymentStub.orderID != response.Order.ID {
		t.Fatalf("expected payment order ID %q, got %q", response.Order.ID, paymentStub.orderID)
	}
}

func TestServiceCheckoutRejectsEmptyCart(t *testing.T) {
	productService := product.NewService(product.NewMemoryRepository())
	cartService := cart.NewService(cart.NewMemoryRepository(), productService)
	paymentStub := &paymentServiceStub{clientSecret: "pi_secret_test"}
	orderService := NewService(NewMemoryRepository(), cartService, productService, paymentStub)

	_, err := orderService.Checkout("user-empty")
	if err == nil {
		t.Fatal("expected empty cart error")
	}

	if got, want := err.Error(), "empty cart"; got != want {
		t.Fatalf("expected error %q, got %q", want, got)
	}

	if paymentStub.calls != 0 {
		t.Fatalf("expected payment service not to be called, got %d calls", paymentStub.calls)
	}
}

func TestServicePayMarksOrderAsPaidDecreasesStockAndClearsCart(t *testing.T) {
	productRepo := product.NewMemoryRepository()
	productService := product.NewService(productRepo)
	createdProduct, err := productService.Create("Mouse", "Wireless mouse", 1500, 5)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	cartRepo := cart.NewMemoryRepository()
	cartService := cart.NewService(cartRepo, productService)
	if _, err := cartService.AddItems("user-1", createdProduct.ID, 2); err != nil {
		t.Fatalf("AddItems() error = %v", err)
	}

	orderRepo := NewMemoryRepository()
	existingOrder := &Order{
		ID:     "order-1",
		UserID: "user-1",
		Status: StatusPending,
		Items: []OrderItem{
			{ProductID: createdProduct.ID, Quantity: 2, Price: createdProduct.Price},
		},
		Total: 3000,
	}
	if err := orderRepo.Create(existingOrder); err != nil {
		t.Fatalf("Create(order) error = %v", err)
	}

	orderService := NewService(orderRepo, cartService, productService, &paymentServiceStub{})

	if err := orderService.Pay("order-1"); err != nil {
		t.Fatalf("Pay() error = %v", err)
	}

	storedOrder, err := orderRepo.FindByID("order-1")
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if storedOrder.Status != StatusPaid {
		t.Fatalf("expected order status %q, got %q", StatusPaid, storedOrder.Status)
	}

	updatedProduct, err := productService.GetByID(createdProduct.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if updatedProduct.Stock != 3 {
		t.Fatalf("expected stock 3, got %d", updatedProduct.Stock)
	}

	userCart, err := cartService.GetCart("user-1")
	if err != nil {
		t.Fatalf("GetCart() error = %v", err)
	}

	if len(userCart.Items) != 0 {
		t.Fatalf("expected cleared cart, got %d items", len(userCart.Items))
	}
}
