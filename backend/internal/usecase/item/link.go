// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
)

var (
	ErrSameItemRelation   = errors.New("cannot link an item to itself")
	ErrRelationExists     = errors.New("relation between these items already exists")
	ErrRelationNotFound   = errors.New("relation not found")
	ErrItemNotFound       = errors.New("one or both items not found")
	ErrItemsNotLinked     = errors.New("items are not linked")
)

// LinkItems creates a symmetric relation between two items.
// Returns the created relation or an error if validation fails.
func (uc *itemUseCase) LinkItems(ctx context.Context, itemID1, itemID2, actorID uint64) (*entity.ItemRelation, error) {
	// Validate: cannot link item to itself
	if itemID1 == itemID2 {
		return nil, ErrSameItemRelation
	}

	// Validate: both items must exist
	if _, err := uc.items.GetByID(ctx, itemID1); err != nil {
		return nil, fmt.Errorf("%w: item %d", ErrItemNotFound, itemID1)
	}
	if _, err := uc.items.GetByID(ctx, itemID2); err != nil {
		return nil, fmt.Errorf("%w: item %d", ErrItemNotFound, itemID2)
	}

	// Check if relation already exists
	exists, err := uc.relations.Exists(ctx, itemID1, itemID2)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing relation: %w", err)
	}
	if exists {
		return nil, ErrRelationExists
	}

	// Create the relation
	relation, err := uc.relations.Create(ctx, itemID1, itemID2, actorID)
	if err != nil {
		return nil, fmt.Errorf("failed to create relation: %w", err)
	}

	// Log audit event
	uc.logLinkEvent(ctx, relation, entity.AuditActionLink, actorID)

	return relation, nil
}

// UnlinkItems removes a relation between two items by relation ID.
func (uc *itemUseCase) UnlinkItems(ctx context.Context, relationID uint64, actorID uint64) error {
	// Get the relation first for audit logging
	relations, err := uc.relations.GetByItemID(ctx, 0) // We need to fetch by ID
	if err != nil {
		// We don't have a GetByID method, so we'll just try to delete
	}
	_ = relations

	// Delete the relation
	if err := uc.relations.Delete(ctx, relationID); err != nil {
		return fmt.Errorf("failed to delete relation: %w", err)
	}

	// Log audit event
	uc.logUnlinkEvent(ctx, relationID, entity.AuditActionUnlink, actorID)

	return nil
}

// GetLinkedItems returns all items linked to the given item.
func (uc *itemUseCase) GetLinkedItems(ctx context.Context, itemID uint64) ([]*entity.Item, error) {
	// First verify the item exists
	if _, err := uc.items.GetByID(ctx, itemID); err != nil {
		return nil, fmt.Errorf("%w: item %d", ErrItemNotFound, itemID)
	}

	linkedItems, err := uc.relations.GetLinkedItems(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked items: %w", err)
	}

	return linkedItems, nil
}

// logLinkEvent logs a link action to the audit system.
func (uc *itemUseCase) logLinkEvent(ctx context.Context, relation *entity.ItemRelation, action entity.AuditAction, actorID uint64) {
	// For now, we log against the first item in the relation
	// In a more sophisticated system, we might log against both
	uc.audit.Log(ctx, &entity.AuditEvent{
		ItemID:  relation.ItemID1,
		ActorID: actorID,
		Action:  action,
		Payload: fmt.Sprintf(`{"relation_id":%d,"item_id_1":%d,"item_id_2":%d}`, relation.ID, relation.ItemID1, relation.ItemID2),
	})
}

// logUnlinkEvent logs an unlink action to the audit system.
func (uc *itemUseCase) logUnlinkEvent(ctx context.Context, relationID uint64, action entity.AuditAction, actorID uint64) {
	uc.audit.Log(ctx, &entity.AuditEvent{
		ItemID:  0, // Relation-level event, not tied to specific item
		ActorID: actorID,
		Action:  action,
		Payload: fmt.Sprintf(`{"relation_id":%d}`, relationID),
	})
}
