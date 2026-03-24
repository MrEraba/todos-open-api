package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Active *bool   `json:"active,omitempty"`
}

// NewUser creates a new user with the given name and email
func NewUser(name string, email string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Active:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
