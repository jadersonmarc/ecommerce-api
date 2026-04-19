package order

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

func (r *PostgresRepository) Create(order *Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&database.OrderModel{
			ID:        order.ID,
			UserID:    order.UserID,
			Total:     order.Total,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.CreatedAt,
		}).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := tx.Create(&database.OrderItemModel{
				OrderID:   order.ID,
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

func (r *PostgresRepository) FindByID(id string) (*Order, error) {
	var model database.OrderModel
	err := r.db.Preload("Items").First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	order := &Order{
		ID:        model.ID,
		UserID:    model.UserID,
		Total:     model.Total,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	for _, item := range model.Items {
		order.Items = append(order.Items, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return order, nil
}

func (r *PostgresRepository) UpdateStatus(id, status string) error {
	result := r.db.Model(&database.OrderModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return database.ErrNotFound
	}

	return nil
}
