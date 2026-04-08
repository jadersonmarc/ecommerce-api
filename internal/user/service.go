package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func generateID() string {
	return uuid.New().String()
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Register(name, email, password string) (*User, error) {
	existing, _ := s.repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		Password:  string(hash),
		Role:      RoleUser,
		UpdatedAt: time.Now(),
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(email, passeword string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passeword))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
