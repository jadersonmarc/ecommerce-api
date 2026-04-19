package database

import (
	"context"

	"gorm.io/gorm"
)

func migrate(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(
		&UserModel{},
		&ProductModel{},
		&CartModel{},
		&CartItemModel{},
		&OrderModel{},
		&OrderItemModel{},
	)
}
