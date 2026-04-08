package user

import "errors"

type MemoryRepository struct {
	users map[string]*User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{users: make(map[string]*User)}
}

func (r *MemoryRepository) Create(user *User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	r.users[user.Email] = user
	return nil
}

func (r *MemoryRepository) FindByEmail(email string) (*User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
