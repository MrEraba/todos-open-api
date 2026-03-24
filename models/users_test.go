package models

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	user := NewUser("John", "john@example.com")

	if user.Name != "John" {
		t.Errorf("expected name 'John', got '%s'", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got '%s'", user.Email)
	}

	if user.Active != false {
		t.Errorf("expected Active to be false, got %v", user.Active)
	}

	if user.ID == "" {
		t.Error("expected ID to be non-empty")
	}
}

func TestUserFields(t *testing.T) {
	user := &User{
		ID:     "test-id",
		Name:   "Test User",
		Email:  "test@example.com",
		Active: true,
	}

	if user.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", user.ID)
	}

	if user.Name != "Test User" {
		t.Errorf("expected Name 'Test User', got '%s'", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got '%s'", user.Email)
	}

	if !user.Active {
		t.Error("expected Active to be true")
	}
}
