package user

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestServiceRegisterHashesPasswordAndSetsDefaults(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	registered, err := service.Register("Marc", "marc@test.com", "123456")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if registered.ID == "" {
		t.Fatal("expected generated user ID")
	}

	if registered.Role != RoleUser {
		t.Fatalf("expected role %q, got %q", RoleUser, registered.Role)
	}

	if registered.Password == "123456" {
		t.Fatal("expected password to be hashed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(registered.Password), []byte("123456")); err != nil {
		t.Fatalf("expected stored hash to match password: %v", err)
	}

	stored, err := repo.FindByEmail("marc@test.com")
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}

	if stored != registered {
		t.Fatal("expected repository to store the registered user")
	}
}

func TestServiceRegisterRejectsDuplicateEmail(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	if _, err := service.Register("Marc", "marc@test.com", "123456"); err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err := service.Register("Marc 2", "marc@test.com", "abcdef")
	if err == nil {
		t.Fatal("expected duplicate email error")
	}

	if got, want := err.Error(), "email already in use"; got != want {
		t.Fatalf("expected error %q, got %q", want, got)
	}
}

func TestServiceLogin(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	registered, err := service.Register("Marc", "marc@test.com", "123456")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	loggedIn, err := service.Login("marc@test.com", "123456")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if loggedIn != registered {
		t.Fatal("expected Login() to return the stored user")
	}
}

func TestServiceLoginRejectsInvalidCredentials(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	if _, err := service.Register("Marc", "marc@test.com", "123456"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{name: "unknown email", email: "missing@test.com", password: "123456"},
		{name: "wrong password", email: "marc@test.com", password: "wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.email, tt.password)
			if err == nil {
				t.Fatal("expected invalid credentials error")
			}

			if got, want := err.Error(), "invalid email or password"; got != want {
				t.Fatalf("expected error %q, got %q", want, got)
			}
		})
	}
}
