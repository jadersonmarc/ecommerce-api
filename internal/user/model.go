package user

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      Role
	UpdatedAt time.Time
}
