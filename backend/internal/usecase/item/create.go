// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// Create validates input, checks uniqueness, and persists a new item.
func (uc *itemUseCase) Create(ctx context.Context, input CreateInput) (*entity.Item, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	exists, err := uc.items.BarcodeExists(ctx, input.Barcode)
	if err != nil {
		return nil, fmt.Errorf("item.Create barcodeExists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("barcode %q: %w", input.Barcode, logger.ErrAlreadyExists)
	}

	invExists, err := uc.items.InventoryNumberExists(ctx, input.InventoryNumber)
	if err != nil {
		return nil, fmt.Errorf("item.Create inventoryNumberExists: %w", err)
	}
	if invExists {
		return nil, fmt.Errorf("inventory_number %q: %w", input.InventoryNumber, logger.ErrAlreadyExists)
	}

	if _, err := uc.categories.GetByID(ctx, input.CategoryID); err != nil {
		return nil, fmt.Errorf("item.Create category: %w", err)
	}

	if _, err := uc.rooms.GetByID(ctx, input.RoomID); err != nil {
		return nil, fmt.Errorf("item.Create room: %w", err)
	}

	item := &entity.Item{
		Barcode:         input.Barcode,
		InventoryNumber: input.InventoryNumber,
		Name:            input.Name,
		CategoryID:      input.CategoryID,
		RoomID:          input.RoomID,
		Description:     input.Description,
		Status:          entity.StatusActive,
		CreatedBy:       input.ActorID,
		LastEditedBy:    input.ActorID,
	}

	if err := uc.items.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("item.Create persist: %w", err)
	}

	item, err = uc.items.GetByID(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("item.Create refetch: %w", err)
	}

	uc.logEvent(ctx, item, entity.ActionCreated, input.ActorID)
	return item, nil
}

// validateCreateInput checks that mandatory fields are present.
func validateCreateInput(input CreateInput) error {
	if input.Barcode == "" {
		return fmt.Errorf("barcode is required: %w", logger.ErrInvalidInput)
	}
	if input.InventoryNumber == "" {
		return fmt.Errorf("inventory_number is required: %w", logger.ErrInvalidInput)
	}
	if input.Name == "" {
		return fmt.Errorf("name is required: %w", logger.ErrInvalidInput)
	}
	if input.CategoryID == 0 {
		return fmt.Errorf("category_id is required: %w", logger.ErrInvalidInput)
	}
	if input.RoomID == 0 {
		return fmt.Errorf("room_id is required: %w", logger.ErrInvalidInput)
	}
	if input.ActorID == 0 {
		return fmt.Errorf("actor_id is required: %w", logger.ErrInvalidInput)
	}
	return nil
}
