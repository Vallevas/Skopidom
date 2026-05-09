// Package repository defines persistence contracts for domain entities.
package repository

import (
	"context"
	"time"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
)

// ItemFilter holds optional parameters for filtering item list queries.
type ItemFilter struct {
	CategoryID *uint64
	RoomID     *uint64
	Status     *entity.ItemStatus
	DateFrom   *time.Time
	DateTo     *time.Time
}

// ItemRepository defines the persistence contract for inventory items.
type ItemRepository interface {
	Create(ctx context.Context, item *entity.Item) error
	GetByID(ctx context.Context, id uint64) (*entity.Item, error)
	GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error)
	List(ctx context.Context, filter ItemFilter) ([]*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) error
	UpdateStatus(ctx context.Context, item *entity.Item) error
	UpdateTxHash(ctx context.Context, id uint64, txHash string) error
	BarcodeExists(ctx context.Context, barcode string) (bool, error)
	InventoryNumberExists(ctx context.Context, inventoryNumber string) (bool, error)
	// MoveToRoom changes the room of an item and records the actor.
	MoveToRoom(ctx context.Context, itemID uint64, roomID uint64, actorID uint64) error
}

// PhotoRepository defines the persistence contract for item photos.
type PhotoRepository interface {
	// Add stores a new photo URL for an item and returns the created record.
	Add(ctx context.Context, photo *entity.ItemPhoto) error
	// ListByItem returns all photos for the given item in upload order.
	ListByItem(ctx context.Context, itemID uint64) ([]*entity.ItemPhoto, error)
	// Delete removes a single photo by its ID, verifying it belongs to itemID.
	Delete(ctx context.Context, photoID uint64, itemID uint64) error
}

// UserRepository defines the persistence contract for system users.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint64) error
	EmailExists(ctx context.Context, email string) (bool, error)
	CountByRole(ctx context.Context, role entity.UserRole) (int, error)
}

// CategoryRepository defines the persistence contract for item categories.
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uint64) (*entity.Category, error)
	GetByName(ctx context.Context, name string) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint64) error
}

// RoomRepository defines the persistence contract for physical rooms.
type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	GetByID(ctx context.Context, id uint64) (*entity.Room, error)
	GetByNameAndBuilding(ctx context.Context, name string, buildingID uint64) (*entity.Room, error)
	ListByBuilding(ctx context.Context, buildingID uint64) ([]*entity.Room, error)
	List(ctx context.Context) ([]*entity.Room, error)
	Update(ctx context.Context, room *entity.Room) error
	Delete(ctx context.Context, id uint64) error
}

// BuildingRepository defines the persistence contract for university buildings.
type BuildingRepository interface {
	Create(ctx context.Context, building *entity.Building) error
	GetByID(ctx context.Context, id uint64) (*entity.Building, error)
	GetByName(ctx context.Context, name string) (*entity.Building, error)
	List(ctx context.Context) ([]*entity.Building, error)
	Update(ctx context.Context, building *entity.Building) error
	Delete(ctx context.Context, id uint64) error
}

// AuditLogger defines the contract for recording item lifecycle events.
type AuditLogger interface {
	Log(ctx context.Context, event *entity.AuditEvent) error
	ListByItem(ctx context.Context, itemID uint64) ([]*entity.AuditEvent, error)
}

// DisposalDocumentRepository defines the persistence contract for disposal documents.
type DisposalDocumentRepository interface {
	// Create stores a new disposal document and returns the created record.
	Create(ctx context.Context, doc *entity.DisposalDocument) error
	// GetByID retrieves a single disposal document by its ID.
	GetByID(ctx context.Context, id uint64) (*entity.DisposalDocument, error)
	// ListByItem returns all disposal documents for the given item.
	ListByItem(ctx context.Context, itemID uint64) ([]*entity.DisposalDocument, error)
	// CountByItem returns the number of disposal documents for the given item.
	CountByItem(ctx context.Context, itemID uint64) (int64, error)
	// Delete removes a disposal document by its ID, verifying it belongs to itemID.
	Delete(ctx context.Context, docID uint64, itemID uint64) error
}

// ItemRelationRepository defines the persistence contract for item relations.
type ItemRelationRepository interface {
	// Create establishes a symmetric relation between two items.
	Create(ctx context.Context, itemID1, itemID2, createdBy uint64) (*entity.ItemRelation, error)
	// GetByItemID returns all relations involving the given item.
	GetByItemID(ctx context.Context, itemID uint64) ([]*entity.ItemRelation, error)
	// GetLinkedItems returns all items linked to the given item.
	GetLinkedItems(ctx context.Context, itemID uint64) ([]*entity.Item, error)
	// Delete removes a relation by its ID.
	Delete(ctx context.Context, relationID uint64) error
	// Exists checks if a relation exists between two items.
	Exists(ctx context.Context, itemID1, itemID2 uint64) (bool, error)
}
