// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	userUC "github.com/Vallevas/Skopidom/internal/usecase/user"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// jwtConfig holds token signing configuration injected at construction time.
type jwtConfig struct {
	secret string
	ttl    time.Duration
}

// AuthHandler handles login and token-related endpoints.
type AuthHandler struct {
	uc  userUC.UseCase
	jwt jwtConfig
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(uc userUC.UseCase, secret string, ttl time.Duration) *AuthHandler {
	return &AuthHandler{uc: uc, jwt: jwtConfig{secret: secret, ttl: ttl}}
}

// UserHandler handles user management endpoints (admin-only).
type UserHandler struct {
	uc userUC.UseCase
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(uc userUC.UseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// ── Auth endpoints ────────────────────────────────────────────────────────────

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  *entity.User `json:"user"`
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	user, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := generateJWT(user, h.jwt)
	if err != nil {
		handleError(w, err)
		return
	}

	respond(w, http.StatusOK, loginResponse{Token: token, User: user})
}

// ── User CRUD endpoints (admin-only) ─────────────────────────────────────────

type registerUserRequest struct {
	FullName string          `json:"full_name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     entity.UserRole `json:"role"`
}

type updateUserRequest struct {
	FullName string          `json:"full_name"`
	Role     entity.UserRole `json:"role"`
}

// Register godoc
// POST /api/v1/users
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	input := userUC.RegisterInput{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	user, err := h.uc.Register(r.Context(), input)
	if err != nil {
		handleError(w, err)
		return
	}

	respond(w, http.StatusCreated, user)
}

// List godoc
// GET /api/v1/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, users)
}

// GetByID godoc
// GET /api/v1/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	user, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, user)
}

// Update godoc
// PATCH /api/v1/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	var req updateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	input := userUC.UpdateInput{
		UserID:   id,
		FullName: req.FullName,
		Role:     req.Role,
	}

	user, err := h.uc.Update(r.Context(), input)
	if err != nil {
		handleError(w, err)
		return
	}

	respond(w, http.StatusOK, user)
}

// Delete godoc
// DELETE /api/v1/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	if err := h.uc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── lookup handlers (categories, rooms) ──────────────────────────────────────

// LookupHandler handles read-only and write endpoints for reference data.
type LookupHandler struct {
	categories lookupUseCase
	rooms      lookupUseCase
}

// lookupUseCase is a minimal interface satisfied by category and room use cases.
type lookupUseCase interface {
	Create(ctx interface{}, name string) (interface{}, error)
}

// ── private helpers ───────────────────────────────────────────────────────────

// generateJWT creates a signed JWT for the given user using the given config.
func generateJWT(user *entity.User, cfg jwtConfig) (string, error) {
	// Import the middleware package's GenerateToken function via an alias
	// to avoid circular imports — actual call is wired in router.go.
	// This placeholder is replaced by the router when constructing handlers.
	_ = user
	_ = cfg
	return "", fmt.Errorf("generateJWT: not wired — use middleware.GenerateToken")
}

// wrapInvalidInput wraps a generic decode error as ErrInvalidInput.
func wrapInvalidInput(err error) error {
	return fmt.Errorf("%s: %w", err.Error(), apperrors.ErrInvalidInput)
}

// urlID is a helper used by handlers that don't use chi directly.
func urlID(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
