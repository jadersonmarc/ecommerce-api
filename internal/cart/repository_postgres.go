package cart

import (
	"errors"

	"github.com/jadersonmarc/ecommerce-api/internal/database"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByUserID(userID string) (*Cart, error) {
	var model database.CartModel
	err := r.db.Preload("Items").First(&model, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	cart := &Cart{
		ID:        model.ID,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	for _, item := range model.Items {
		cart.Items = append(cart.Items, CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return cart, nil
}

func (r *PostgresRepository) Save(cart *Cart) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		model := database.CartModel{
			ID:        cart.ID,
			UserID:    cart.UserID,
			CreatedAt: cart.CreatedAt,
			UpdatedAt: cart.UpdatedAt,
		}

		var existing database.CartModel
		err := tx.Select("id", "user_id", "created_at").First(&existing, "user_id = ?", cart.UserID).Error
		if err == nil {
			model.ID = existing.ID
			model.CreatedAt = existing.CreatedAt
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Save(&model).Error; err != nil {
			return err
		}

		cart.ID = model.ID
		cart.CreatedAt = model.CreatedAt

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&database.CartItemModel{}).Error; err != nil {
			return err
		}

		for _, item := range cart.Items {
			if err := tx.Create(&database.CartItemModel{
				CartID:    cart.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
