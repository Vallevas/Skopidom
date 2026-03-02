// Package repository defines persistence contracts for domain entities.
package repository

import (
	"context"
	"time"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
)

// ItemFilter holds optional parameters for filtering item list queries.
// Nil pointer fields are ignored (not applied to the query).
type ItemFilter struct {
	CategoryID *uint64
	RoomID     *uint64
	Status     *entity.ItemStatus
	DateFrom   *time.Time
	DateTo     *time.Time
}

// ItemRepository defines the persistence contract for inventory items.
type ItemRepository interface {
	// Create persists a new item and sets its generated ID and timestamps.
	Create(ctx context.Context, item *entity.Item) error

	// GetByID returns the item with the given ID or ErrNotFound.
	GetByID(ctx context.Context, id uint64) (*entity.Item, error)

	// GetByBarcode returns the item matching the barcode or ErrNotFound.
	GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error)

	// List returns items matching the provided filter.
	List(ctx context.Context, filter ItemFilter) ([]*entity.Item, error)

	// Update persists changes to Description, PhotoURL, and UpdatedAt.
	Update(ctx context.Context, item *entity.Item) error

	// UpdateStatus persists a status change (e.g. disposed).
	UpdateStatus(ctx context.Context, item *entity.Item) error

	// UpdateTxHash stores the blockchain transaction hash for an item.
	UpdateTxHash(ctx context.Context, id uint64, txHash string) error

	// BarcodeExists reports whether a barcode is already registered.
	BarcodeExists(ctx context.Context, barcode string) (bool, error)
}

// UserRepository defines the persistence contract for system users.
type UserRepository interface {
	// Create persists a new user and sets its generated ID and timestamps.
	Create(ctx context.Context, user *entity.User) error

	// GetByID returns the user with the given ID or ErrNotFound.
	GetByID(ctx context.Context, id uint64) (*entity.User, error)

	// GetByEmail returns the user with the given email or ErrNotFound.
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// List returns all users.
	List(ctx context.Context) ([]*entity.User, error)

	// Update persists changes to FullName, Role, and UpdatedAt.
	Update(ctx context.Context, user *entity.User) error

	// Delete permanently removes a user record.
	Delete(ctx context.Context, id uint64) error

	// EmailExists reports whether an email address is already registered.
	EmailExists(ctx context.Context, email string) (bool, error)
}

// CategoryRepository defines the persistence contract for item categories.
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uint64) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint64) error
}

// RoomRepository defines the persistence contract for physical rooms.
type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	GetByID(ctx context.Context, id uint64) (*entity.Room, error)
	// ListByBuilding returns all rooms belonging to the given building.
	ListByBuilding(ctx context.Context, buildingID uint64) ([]*entity.Room, error)
	List(ctx context.Context) ([]*entity.Room, error)
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, id uint64) error
}

// BuildingRepository defines the persistence contract for university buildings.
type BuildingRepository interface {
	Create(ctx context.Context, building *entity.Building) error
	GetByID(ctx context.Context, id uint64) (*entity.Building, error)
	List(ctx context.Context) ([]*entity.Building, error)
	Update(ctx context.Context, building *entity.Building) error
	Delete(ctx context.Context, id uint64) error
}
