// Package middleware provides HTTP middleware components for the inventory API.
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

const (
	// ContextKeyUserID is the context key for the authenticated user's ID.
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyUserRole is the context key for the authenticated user's role.
	ContextKeyUserRole contextKey = "user_role"
)

// Claims holds the JWT payload fields.
type Claims struct {
	UserID uint64          `json:"uid"`
	Role   entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(userID uint64, role entity.UserRole, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Auth returns middleware that validates the Bearer JWT in the Authorization header.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractBearerToken(r)
			if tokenString == "" {
				writeError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(
				tokenString, claims,
				func(_ *jwt.Token) (any, error) { return []byte(secret), nil },
			)
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized)
				return
			}

			// Inject user info into context for downstream handlers.
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyUserRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin returns middleware that allows only admin-role users through.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(ContextKeyUserRole).(entity.UserRole)
		if role != entity.RoleAdmin {
			writeError(w, http.StatusForbidden, apperrors.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserIDFromCtx extracts the authenticated user ID from a request context.
func UserIDFromCtx(ctx context.Context) uint64 {
	id, _ := ctx.Value(ContextKeyUserID).(uint64)
	return id
}

// UserRoleFromCtx extracts the authenticated user role from a request context.
func UserRoleFromCtx(ctx context.Context) entity.UserRole {
	role, _ := ctx.Value(ContextKeyUserRole).(entity.UserRole)
	return role
}

// extractBearerToken parses the "Bearer <token>" Authorization header value.
func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}

// writeError writes a JSON error response with the given HTTP status code.
func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + err.Error() + `"}`))
}
