// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
)

// UseCase defines all business operations on inventory items.
type UseCase interface {
	// Create registers a new item from a barcode scan or manual entry.
	Create(ctx context.Context, input CreateInput) (*entity.Item, error)

	// Update modifies the mutable fields of an existing item.
	Update(ctx context.Context, input UpdateInput) (*entity.Item, error)

	// Dispose transitions an item to the disposed lifecycle state.
	Dispose(ctx context.Context, itemID uint64, actorID uint64) error

	// GetByID returns a single item with its relations.
	GetByID(ctx context.Context, id uint64) (*entity.Item, error)

	// GetByBarcode returns a single item by barcode value.
	GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error)

	// List returns items matching the given filter.
	List(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, error)

	// ListAuditEvents returns the full audit history for the given item.
	ListAuditEvents(ctx context.Context, itemID uint64) ([]*entity.AuditEvent, error)
}

// CreateInput holds the data required to register a new inventory item.
type CreateInput struct {
	Barcode     string
	Name        string
	CategoryID  uint64
	RoomID      uint64
	Description string
	PhotoURL    string
	// ActorID is the ID of the user performing the action.
	ActorID uint64
}

// UpdateInput holds the data allowed to be changed after creation.
type UpdateInput struct {
	// ItemID identifies the record to update.
	ItemID      uint64
	Description string
	PhotoURL    string
	// ActorID is the ID of the user performing the action.
	ActorID uint64
}

// itemUseCase is the concrete implementation of UseCase.
type itemUseCase struct {
	items      repository.ItemRepository
	categories repository.CategoryRepository
	rooms      repository.RoomRepository
	audit      repository.AuditLogger
}

// New constructs an itemUseCase with all required repository dependencies.
func New(
	items repository.ItemRepository,
	categories repository.CategoryRepository,
	rooms repository.RoomRepository,
	audit repository.AuditLogger,
) UseCase {
	return &itemUseCase{
		items:      items,
		categories: categories,
		rooms:      rooms,
		audit:      audit,
	}
}

// ListAuditEvents returns the full audit history for the given item.
func (uc *itemUseCase) ListAuditEvents(
	ctx context.Context,
	itemID uint64,
) ([]*entity.AuditEvent, error) {
	return uc.audit.ListByItem(ctx, itemID)
}

// logEvent is a helper that builds and persists an AuditEvent.
// It never fails the caller — errors are logged internally.
func (uc *itemUseCase) logEvent(
	ctx context.Context,
	item *entity.Item,
	action entity.AuditAction,
	actorID uint64,
) {
	payload, err := json.Marshal(item)
	if err != nil {
		slog.Error("audit: failed to marshal item snapshot",
			"item_id", item.ID, "err", err)
		return
	}

	_ = uc.audit.Log(ctx, &entity.AuditEvent{
		ItemID:  item.ID,
		ActorID: actorID,
		Action:  action,
		Payload: string(payload),
	})
}

