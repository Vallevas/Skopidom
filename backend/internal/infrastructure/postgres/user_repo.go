// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// UserRepo implements repository.UserRepository using sqlc-generated queries.
type UserRepo struct {
	queries *db.Queries
}

// NewUserRepo constructs a UserRepo backed by the given connection pool.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Create inserts a new user and populates the generated ID and timestamps.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
	})
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}

	user.ID = uint64(row.ID)
	user.CreatedAt = row.CreatedAt
	user.UpdatedAt = row.UpdatedAt
	return nil
}

// GetByID returns the user matching the given ID or ErrNotFound.
func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	row, err := r.queries.GetUserByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByID: %w", err)
	}
	return mapUser(row), nil
}

// GetByEmail returns the user matching the given email or ErrNotFound.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}
	return mapUser(row), nil
}

// List returns all registered users ordered by creation date.
func (r *UserRepo) List(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.List: %w", err)
	}

	users := make([]*entity.User, len(rows))
	for i, row := range rows {
		users[i] = mapUser(row)
	}
	return users, nil
}

// Update persists changes to FullName and Role fields.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	updatedAt, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		FullName: user.FullName,
		Role:     string(user.Role),
		ID:       int64(user.ID),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	user.UpdatedAt = updatedAt
	return nil
}

// Delete permanently removes a user record from the database.
func (r *UserRepo) Delete(ctx context.Context, id uint64) error {
	if err := r.queries.DeleteUser(ctx, int64(id)); err != nil {
		return fmt.Errorf("UserRepo.Delete: %w", err)
	}
	return nil
}

// EmailExists reports whether the given email address is already registered.
func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.EmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("UserRepo.EmailExists: %w", err)
	}
	return exists, nil
}

// CountByRole returns the number of users assigned the given role.
func (r *UserRepo) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	count, err := r.queries.CountUsersByRole(ctx, string(role))
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CountByRole: %w", err)
	}
	return int(count), nil
}

// ── mapping ───────────────────────────────────────────────────────────────────

// mapUser converts a sqlc-generated db.User to a domain entity.
// GetUserByID, GetUserByEmail and ListUsers all return db.User — one mapper handles all.
func mapUser(row db.User) *entity.User {
	return &entity.User{
		ID:           uint64(row.ID),
		FullName:     row.FullName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Role:         entity.UserRole(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

