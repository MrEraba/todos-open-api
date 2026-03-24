package store

import (
	"errors"
	"sync"
	"time"

	"github.com/MrEraba/todos-open-api/models"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicate    = errors.New("resource already exists")
)

// UserStore provides in-memory storage for users
type UserStore struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// NewUserStore creates a new in-memory user store
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*models.User),
	}
}

// Create adds a new user to the store
func (s *UserStore) Create(req *models.CreateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate email
	for _, u := range s.users {
		if u.Email == req.Email {
			return nil, ErrDuplicate
		}
	}

	user := &models.User{
		ID:        generateID(),
		Name:      req.Name,
		Email:     req.Email,
		Active:    false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	s.users[user.ID] = user
	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserStore) GetByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	return user, nil
}

// GetAll retrieves all users
func (s *UserStore) GetAll() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// Update updates an existing user
func (s *UserStore) Update(id string, req *models.UpdateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Check for duplicate email if email is being changed
	if req.Email != nil && *req.Email != user.Email {
		for _, u := range s.users {
			if u.Email == *req.Email && u.ID != id {
				return nil, ErrDuplicate
			}
		}
		user.Email = *req.Email
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	user.UpdatedAt = time.Now().UTC()
	return user, nil
}

// Delete removes a user from the store
func (s *UserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return ErrNotFound
	}

	delete(s.users, id)
	return nil
}

// Count returns the number of users in the store
func (s *UserStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}

// generateID creates a unique ID using UUID
func generateID() string {
	return uuid.New().String()
}
