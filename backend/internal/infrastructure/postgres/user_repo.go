// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepo implements repository.UserRepository using PostgreSQL.
type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo constructs a UserRepo backed by the given connection pool.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create inserts a new user and populates the generated ID and timestamps.
func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	return nil
}

// GetByID returns the user matching the given ID or ErrNotFound.
func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &entity.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email,
		&user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByID: %w", err)
	}
	return user, nil
}

// GetByEmail returns the user matching the given email or ErrNotFound.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &entity.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.FullName, &user.Email,
		&user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}
	return user, nil
}

// List returns all registered users ordered by creation date.
func (r *UserRepo) List(ctx context.Context) ([]*entity.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, created_at, updated_at
		FROM users
		ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.List: %w", err)
	}
	defer rows.Close()

	users := make([]*entity.User, 0)
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(
			&user.ID, &user.FullName, &user.Email,
			&user.PasswordHash, &user.Role,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("UserRepo.List scan: %w", err)
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Update persists changes to FullName and Role fields.
func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET full_name = $1, role = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		user.FullName,
		user.Role,
		user.ID,
	).Scan(&user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", err)
	}
	return nil
}

// Delete permanently removes a user record from the database.
func (r *UserRepo) Delete(ctx context.Context, id uint64) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("UserRepo.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return logger.ErrNotFound
	}
	return nil
}

// EmailExists reports whether the given email address is already registered.
func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("UserRepo.EmailExists: %w", err)
	}
	return exists, nil
}

// CountByRole returns the number of users assigned the given role.
func (r *UserRepo) CountByRole(ctx context.Context, role entity.UserRole) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE role = $1`, role,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CountByRole: %w", err)
	}
	return count, nil
}

