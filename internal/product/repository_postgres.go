package product

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

func (r *PostgresRepository) Create(product *Product) error {
	return r.db.Create(&database.ProductModel{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.createdAt,
	}).Error
}

func (r *PostgresRepository) FindAll() ([]*Product, error) {
	var models []database.ProductModel
	if err := r.db.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	var products []*Product
	for _, model := range models {
		products = append(products, &Product{
			ID:          model.ID,
			Name:        model.Name,
			Description: model.Description,
			Price:       model.Price,
			Stock:       model.Stock,
			createdAt:   model.CreatedAt,
		})
	}

	return products, nil
}

func (r *PostgresRepository) FindByID(id string) (*Product, error) {
	var model database.ProductModel
	err := r.db.First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       model.Price,
		Stock:       model.Stock,
		createdAt:   model.CreatedAt,
	}, nil
}

func (r *PostgresRepository) DecreaseStock(productID string, quantity int) error {
	result := r.db.Model(&database.ProductModel{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	var existing database.ProductModel
	err := r.db.Select("id").First(&existing, "id = ?", productID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.ErrNotFound
	}
	if err != nil {
		return err
	}

	return errors.New("insufficient stock")
}
