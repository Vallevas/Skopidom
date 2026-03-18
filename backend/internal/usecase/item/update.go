// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// Update changes the mutable fields of an existing active item.
func (uc *itemUseCase) Update(ctx context.Context, input UpdateInput) (*entity.Item, error) {
	if input.ItemID == 0 || input.ActorID == 0 {
		return nil, fmt.Errorf("item_id and actor_id are required: %w",
			logger.ErrInvalidInput)
	}

	item, err := uc.items.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, fmt.Errorf("item.Update fetch: %w", err)
	}

	// Guard: disposed items must not be edited.
	if !item.IsMutable() {
		return nil, fmt.Errorf("item %d: %w", input.ItemID, logger.ErrDisposed)
	}

	item.Description = input.Description
	item.PhotoURL = input.PhotoURL
	item.LastEditedBy = input.ActorID

	if err := uc.items.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("item.Update persist: %w", err)
	}

	return uc.items.GetByID(ctx, item.ID)
}

// Dispose transitions an item to the disposed state.
// Only admins should be able to call this — enforce at handler level.
func (uc *itemUseCase) Dispose(ctx context.Context, itemID uint64, actorID uint64) error {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("item.Dispose fetch: %w", err)
	}

	if !item.IsActive() {
		return fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}

	item.Dispose(actorID)

	if err := uc.items.UpdateStatus(ctx, item); err != nil {
		return fmt.Errorf("item.Dispose persist: %w", err)
	}
	return nil
}

// GetByID returns a single item with all relations populated.
func (uc *itemUseCase) GetByID(
	ctx context.Context,
	id uint64,
) (*entity.Item, error) {
	item, err := uc.items.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("item.GetByID: %w", err)
	}
	return item, nil
}

// GetByBarcode returns the item matching the given barcode string.
func (uc *itemUseCase) GetByBarcode(
	ctx context.Context,
	barcode string,
) (*entity.Item, error) {
	if barcode == "" {
		return nil, fmt.Errorf("barcode is required: %w", logger.ErrInvalidInput)
	}
	item, err := uc.items.GetByBarcode(ctx, barcode)
	if err != nil {
		return nil, fmt.Errorf("item.GetByBarcode: %w", err)
	}
	return item, nil
}

// List returns items matching the provided filter.
func (uc *itemUseCase) List(
	ctx context.Context,
	filter repository.ItemFilter,
) ([]*entity.Item, error) {
	items, err := uc.items.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("item.List: %w", err)
	}
	return items, nil
}
