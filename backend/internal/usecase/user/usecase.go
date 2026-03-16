// Package user contains the business logic for user account management.
package user

import (
	"context"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// UseCase defines all business operations on user accounts.
type UseCase interface {
	// Register creates a new user account (admin-only operation).
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)

	// Login authenticates credentials and returns the user on success.
	Login(ctx context.Context, email, password string) (*entity.User, error)

	// GetByID returns the user with the given ID.
	GetByID(ctx context.Context, id uint64) (*entity.User, error)

	// List returns all registered users.
	List(ctx context.Context) ([]*entity.User, error)

	// Update changes the FullName and/or Role of a user.
	Update(ctx context.Context, input UpdateInput) (*entity.User, error)

	// Delete permanently removes a user account.
	// actorID is the ID of the admin performing the deletion.
	Delete(ctx context.Context, id uint64, actorID uint64) error
}

// RegisterInput holds data required to create a new user account.
type RegisterInput struct {
	FullName string
	Email    string
	Password string
	Role     entity.UserRole
}

// UpdateInput holds the fields that can be changed on an existing user.
type UpdateInput struct {
	UserID  uint64
	ActorID uint64
	// FullName and Role are optional — empty value means no change.
	FullName string
	Role     entity.UserRole
}

// userUseCase is the concrete implementation of UseCase.
type userUseCase struct {
	users repository.UserRepository
}

// New constructs a userUseCase with the required repository dependency.
func New(users repository.UserRepository) UseCase {
	return &userUseCase{users: users}
}

// Register validates input, hashes the password, and persists the user.
func (uc *userUseCase) Register(
	ctx context.Context,
	input RegisterInput,
) (*entity.User, error) {
	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	emailTaken, err := uc.users.EmailExists(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("user.Register emailExists: %w", err)
	}
	if emailTaken {
		return nil, fmt.Errorf("email %q: %w", input.Email, logger.ErrAlreadyExists)
	}

	// Use bcrypt cost 12 — sufficient for a university-scale system.
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("user.Register bcrypt: %w", err)
	}

	role := input.Role
	if role == "" {
		role = entity.RoleEditor
	}

	user := &entity.User{
		FullName:     input.FullName,
		Email:        input.Email,
		PasswordHash: string(hashBytes),
		Role:         role,
	}

	if err := uc.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("user.Register persist: %w", err)
	}

	// Clear hash before returning for safety.
	user.PasswordHash = ""
	return user, nil
}

// Login verifies credentials and returns the authenticated user.
func (uc *userUseCase) Login(
	ctx context.Context,
	email, password string,
) (*entity.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required: %w",
			logger.ErrInvalidInput)
	}

	user, err := uc.users.GetByEmail(ctx, email)
	if err != nil {
		// Return generic error to avoid leaking whether the email exists.
		return nil, logger.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return nil, logger.ErrUnauthorized
	}

	user.PasswordHash = ""
	return user, nil
}

// GetByID returns the user matching the given ID.
func (uc *userUseCase) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	user, err := uc.users.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	user.PasswordHash = ""
	return user, nil
}

// List returns all users without password hashes.
func (uc *userUseCase) List(ctx context.Context) ([]*entity.User, error) {
	users, err := uc.users.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("user.List: %w", err)
	}
	for _, usr := range users {
		usr.PasswordHash = ""
	}
	return users, nil
}

// Update applies changes to FullName and Role.
func (uc *userUseCase) Update(
	ctx context.Context,
	input UpdateInput,
) (*entity.User, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user_id is required: %w", logger.ErrInvalidInput)
	}

	user, err := uc.users.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user.Update fetch: %w", err)
	}

	// Prevent an admin from downgrading their own role when they are the
	// last admin — results in a system with no way to manage users.
	if input.UserID == input.ActorID &&
		user.Role == entity.RoleAdmin &&
		input.Role != "" &&
		input.Role != entity.RoleAdmin {
		count, err := uc.users.CountByRole(ctx, entity.RoleAdmin)
		if err != nil {
			return nil, fmt.Errorf("user.Update countAdmins: %w", err)
		}
		if count <= 1 {
			return nil, fmt.Errorf(
				"cannot downgrade the last admin account: %w", logger.ErrForbidden,
			)
		}
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	if err := uc.users.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("user.Update persist: %w", err)
	}

	user.PasswordHash = ""
	return user, nil
}

// Delete removes a user account permanently.
// Guards:
//   - An admin cannot delete their own account.
//   - The last admin account cannot be deleted — this would lock out the system.
func (uc *userUseCase) Delete(ctx context.Context, id uint64, actorID uint64) error {
	if id == actorID {
		return fmt.Errorf("cannot delete own account: %w", logger.ErrForbidden)
	}

	target, err := uc.users.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user.Delete fetch: %w", err)
	}

	if target.Role == entity.RoleAdmin {
		count, err := uc.users.CountByRole(ctx, entity.RoleAdmin)
		if err != nil {
			return fmt.Errorf("user.Delete countAdmins: %w", err)
		}
		if count <= 1 {
			return fmt.Errorf(
				"cannot delete the last admin account: %w", logger.ErrForbidden,
			)
		}
	}

	if err := uc.users.Delete(ctx, id); err != nil {
		return fmt.Errorf("user.Delete: %w", err)
	}
	return nil
}

// validateRegisterInput checks that all required fields are non-empty.
func validateRegisterInput(input RegisterInput) error {
	if input.FullName == "" {
		return fmt.Errorf("full_name is required: %w", logger.ErrInvalidInput)
	}
	if input.Email == "" {
		return fmt.Errorf("email is required: %w", logger.ErrInvalidInput)
	}
	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters: %w",
			logger.ErrInvalidInput)
	}
	return nil
}
