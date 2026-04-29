// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// movePayload is the JSON structure recorded in the audit log for move events.
type movePayload struct {
	FromRoomID       uint64 `json:"from_room_id"`
	FromRoomName     string `json:"from_room_name"`
	FromBuildingName string `json:"from_building_name"`
	ToRoomID         uint64 `json:"to_room_id"`
	ToRoomName       string `json:"to_room_name"`
	ToBuildingName   string `json:"to_building_name"`
}

// Update changes the mutable fields of an existing item.
func (uc *itemUseCase) Update(ctx context.Context, input UpdateInput) (*entity.Item, error) {
	if input.ItemID == 0 || input.ActorID == 0 {
		return nil, fmt.Errorf("item_id and actor_id are required: %w", logger.ErrInvalidInput)
	}

	item, err := uc.items.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, fmt.Errorf("item.Update fetch: %w", err)
	}
	if !item.IsMutable() {
		return nil, fmt.Errorf("item %d: %w", input.ItemID, logger.ErrDisposed)
	}

	item.Description = input.Description
	item.LastEditedBy = input.ActorID

	if err := uc.items.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("item.Update persist: %w", err)
	}

	item, err = uc.items.GetByID(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("item.Update refetch: %w", err)
	}

	uc.logEvent(ctx, item, entity.ActionUpdated, input.ActorID)
	return item, nil
}

// ToggleRepair switches an active item to in_repair and vice versa.
// active    → in_repair  : logs ActionSentToRepair
// in_repair → active     : logs ActionReturnedFromRepair
// disposed  → error
func (uc *itemUseCase) ToggleRepair(
	ctx context.Context,
	itemID uint64,
	actorID uint64,
) (*entity.Item, error) {
	if itemID == 0 || actorID == 0 {
		return nil, fmt.Errorf("item_id and actor_id are required: %w", logger.ErrInvalidInput)
	}

	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("item.ToggleRepair fetch: %w", err)
	}
	if !item.IsMutable() {
		return nil, fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}

	var newStatus entity.ItemStatus
	var action entity.AuditAction

	switch item.Status {
	case entity.StatusActive:
		newStatus = entity.StatusInRepair
		action = entity.ActionSentToRepair
	case entity.StatusInRepair:
		newStatus = entity.StatusActive
		action = entity.ActionReturnedFromRepair
	default:
		return nil, fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}

	item.Status = newStatus
	item.LastEditedBy = actorID

	if err := uc.items.UpdateStatus(ctx, item); err != nil {
		return nil, fmt.Errorf("item.ToggleRepair persist: %w", err)
	}

	item, err = uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("item.ToggleRepair refetch: %w", err)
	}

	uc.logEvent(ctx, item, action, actorID)
	return item, nil
}

// MoveToRoom moves an item to a different room and records from/to details.
func (uc *itemUseCase) MoveToRoom(
	ctx context.Context,
	itemID uint64,
	roomID uint64,
	actorID uint64,
) (*entity.Item, error) {
	if itemID == 0 || roomID == 0 || actorID == 0 {
		return nil, fmt.Errorf("item_id, room_id and actor_id are required: %w",
			logger.ErrInvalidInput)
	}

	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("item.MoveToRoom fetch: %w", err)
	}
	if !item.IsMutable() {
		return nil, fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}
	if item.RoomID == roomID {
		return item, nil
	}

	targetRoom, err := uc.rooms.GetByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("item.MoveToRoom room: %w", err)
	}

	payload := movePayload{
		FromRoomID:       item.RoomID,
		FromRoomName:     item.Room.Name,
		FromBuildingName: item.Room.Building.Name,
		ToRoomID:         targetRoom.ID,
		ToRoomName:       targetRoom.Name,
		ToBuildingName:   targetRoom.Building.Name,
	}

	if err := uc.items.MoveToRoom(ctx, itemID, roomID, actorID); err != nil {
		return nil, fmt.Errorf("item.MoveToRoom persist: %w", err)
	}

	item, err = uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("item.MoveToRoom refetch: %w", err)
	}

	payloadBytes, _ := json.Marshal(payload)
	_ = uc.audit.Log(ctx, &entity.AuditEvent{
		ItemID:  item.ID,
		ActorID: actorID,
		Action:  entity.ActionMoved,
		Payload: string(payloadBytes),
	})

	return item, nil
}

// GetByID returns a single item with all relations populated.
func (uc *itemUseCase) GetByID(ctx context.Context, id uint64) (*entity.Item, error) {
	item, err := uc.items.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("item.GetByID: %w", err)
	}
	return item, nil
}

// GetByBarcode returns the item matching the given barcode string.
func (uc *itemUseCase) GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error) {
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

// AddPhoto stores a new photo for an item.
func (uc *itemUseCase) AddPhoto(
	ctx context.Context,
	itemID uint64,
	url string,
	actorID uint64,
) (*entity.ItemPhoto, error) {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("item.AddPhoto fetch: %w", err)
	}
	if !item.IsMutable() {
		return nil, fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}
	photo := &entity.ItemPhoto{ItemID: itemID, URL: url}
	if err := uc.photos.Add(ctx, photo); err != nil {
		return nil, fmt.Errorf("item.AddPhoto persist: %w", err)
	}
	return photo, nil
}

// DeletePhoto removes a photo from an item.
func (uc *itemUseCase) DeletePhoto(
	ctx context.Context,
	itemID uint64,
	photoID uint64,
	actorID uint64,
) error {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("item.DeletePhoto fetch: %w", err)
	}
	if !item.IsMutable() {
		return fmt.Errorf("item %d: %w", itemID, logger.ErrDisposed)
	}
	if err := uc.photos.Delete(ctx, photoID, itemID); err != nil {
		return fmt.Errorf("item.DeletePhoto: %w", err)
	}
	return nil
}
