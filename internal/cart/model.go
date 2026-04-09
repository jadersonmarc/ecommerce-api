package cart

import "time"

type Cart struct {
	ID        string
	UserID    string
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ProductID string
	Quantity  int
	Price     int64
}
