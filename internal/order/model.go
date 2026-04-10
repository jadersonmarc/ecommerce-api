package order

import "time"

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Total     int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     int64
}

const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusFailed  = "failed"
)
