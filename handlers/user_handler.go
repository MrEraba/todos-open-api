package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/MrEraba/todos-open-api/models"
	"github.com/MrEraba/todos-open-api/store"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	store *store.UserStore
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s *store.UserStore) *UserHandler {
	return &UserHandler{store: s}
}

// List godoc
// @Summary List all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := h.store.GetAll()
	respondJSON(w, http.StatusOK, users)
}

// Get godoc
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := extractID(r)
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "User ID is required")
		return
	}

	user, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON in request body")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "validation_error", "Name is required")
		return
	}
	if req.Email == "" {
		respondError(w, http.StatusBadRequest, "validation_error", "Email is required")
		return
	}

	user, err := h.store.Create(&req)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			respondError(w, http.StatusConflict, "duplicate", "A user with this email already exists")
			return
		}
		slog.Error("failed to create user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// Update godoc
// @Summary Update a user
// @Description Update an existing user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "Fields to update"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r)
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "User ID is required")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON in request body")
		return
	}
	defer r.Body.Close()

	user, err := h.store.Update(id, &req)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		if errors.Is(err, store.ErrDuplicate) {
			respondError(w, http.StatusConflict, "duplicate", "A user with this email already exists")
			return
		}
		slog.Error("failed to update user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r)
	if id == "" {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "User ID is required")
		return
	}

	err := h.store.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			respondError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}
		slog.Error("failed to delete user", "error", err)
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("failed to encode JSON response", "error", err)
		}
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, code, message string) {
	errResp := models.ErrorResponse{
		Error:   code,
		Message: message,
	}
	respondJSON(w, status, errResp)
}

// extractID extracts the user ID from the URL path
func extractID(r *http.Request) string {
	// Expected path: /users/{id}
	path := r.URL.Path
	id := ""
	if len(path) > 7 { // len("/users/") = 7
		id = path[7:]
	}
	return id
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "todos-api",
	}
	respondJSON(w, http.StatusOK, response)
}

// Status godoc
// @Summary Get API status
// @Description Get the current status of the API including user count
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /status [get]
func Status(w http.ResponseWriter, r *http.Request, userCount int) {
	response := map[string]interface{}{
		"status":      "running",
		"service":     "todos-api",
		"user_count":  userCount,
		"server_info": "Go HTTP Server",
	}
	respondJSON(w, http.StatusOK, response)
}
