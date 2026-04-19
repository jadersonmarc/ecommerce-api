package user

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

func (r *PostgresRepository) Create(user *User) error {
	return r.db.Create(&database.UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      string(user.Role),
		UpdatedAt: user.UpdatedAt,
	}).Error
}

func (r *PostgresRepository) FindByEmail(email string) (*User, error) {
	var model database.UserModel

	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Password:  model.Password,
		Role:      Role(model.Role),
		UpdatedAt: model.UpdatedAt,
	}, nil
}
