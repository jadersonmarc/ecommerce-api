package database

import "time"

type UserModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name;not null"`
	Email     string    `gorm:"column:email;not null;uniqueIndex"`
	Password  string    `gorm:"column:password;not null"`
	Role      string    `gorm:"column:role;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (UserModel) TableName() string { return "users" }

type ProductModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name;not null;index:idx_products_name"`
	Description string `gorm:"column:description;not null"`
	Price       int64  `gorm:"column:price;not null"`
	Stock       int    `gorm:"column:stock;not null"`
	CreatedAt   int64  `gorm:"column:created_at;not null"`
}

func (ProductModel) TableName() string { return "products" }

type CartModel struct {
	ID        string          `gorm:"column:id;primaryKey"`
	UserID    string          `gorm:"column:user_id;not null;uniqueIndex"`
	CreatedAt time.Time       `gorm:"column:created_at;not null"`
	UpdatedAt time.Time       `gorm:"column:updated_at;not null"`
	Items     []CartItemModel `gorm:"foreignKey:CartID;references:ID;constraint:OnDelete:CASCADE"`
}

func (CartModel) TableName() string { return "carts" }

type CartItemModel struct {
	CartID    string `gorm:"column:cart_id;primaryKey"`
	ProductID string `gorm:"column:product_id;primaryKey"`
	Quantity  int    `gorm:"column:quantity;not null"`
	Price     int64  `gorm:"column:price;not null"`
}

func (CartItemModel) TableName() string { return "cart_items" }

type OrderModel struct {
	ID        string           `gorm:"column:id;primaryKey"`
	UserID    string           `gorm:"column:user_id;not null;index:idx_orders_user_id"`
	Total     int64            `gorm:"column:total;not null"`
	Status    string           `gorm:"column:status;not null"`
	CreatedAt time.Time        `gorm:"column:created_at;not null"`
	UpdatedAt time.Time        `gorm:"column:updated_at;not null"`
	Items     []OrderItemModel `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE"`
}

func (OrderModel) TableName() string { return "orders" }

type OrderItemModel struct {
	OrderID   string `gorm:"column:order_id;primaryKey"`
	ProductID string `gorm:"column:product_id;primaryKey"`
	Quantity  int    `gorm:"column:quantity;not null"`
	Price     int64  `gorm:"column:price;not null"`
}

func (OrderItemModel) TableName() string { return "order_items" }
