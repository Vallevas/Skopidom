// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// ItemRelationRepo implements repository.ItemRelationRepository using sqlc-generated queries.
type ItemRelationRepo struct {
	queries *db.Queries
}

// NewItemRelationRepo constructs a new ItemRelationRepo backed by the given connection pool.
func NewItemRelationRepo(pool *pgxpool.Pool) *ItemRelationRepo {
	return &ItemRelationRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Create establishes a symmetric relation between two items.
// Ensures itemID1 < itemID2 for consistent storage.
func (r *ItemRelationRepo) Create(ctx context.Context, itemID1, itemID2, createdBy uint64) (*entity.ItemRelation, error) {
	// Enforce ordering: smaller ID first
	if itemID1 > itemID2 {
		itemID1, itemID2 = itemID2, itemID1
	}

	row, err := r.queries.CreateItemRelation(ctx, db.CreateItemRelationParams{
		ItemID1:   int64(itemID1),
		ItemID2:   int64(itemID2),
		CreatedBy: int64(createdBy),
	})
	if err != nil {
		return nil, fmt.Errorf("ItemRelationRepo.Create: %w", err)
	}

	return &entity.ItemRelation{
		ID:        uint64(row.ID),
		ItemID1:   uint64(row.ItemID1),
		ItemID2:   uint64(row.ItemID2),
		CreatedAt: row.CreatedAt,
		CreatedBy: uint64(row.CreatedBy),
	}, nil
}

// GetByItemID returns all relations involving the given item.
func (r *ItemRelationRepo) GetByItemID(ctx context.Context, itemID uint64) ([]*entity.ItemRelation, error) {
	rows, err := r.queries.GetRelationsByItemID(ctx, int64(itemID))
	if err != nil {
		return nil, fmt.Errorf("ItemRelationRepo.GetByItemID: %w", err)
	}

	relations := make([]*entity.ItemRelation, len(rows))
	for i, row := range rows {
		relations[i] = &entity.ItemRelation{
			ID:        uint64(row.ID),
			ItemID1:   uint64(row.ItemID1),
			ItemID2:   uint64(row.ItemID2),
			CreatedAt: row.CreatedAt,
			CreatedBy: uint64(row.CreatedBy),
		}
	}
	return relations, nil
}

// GetLinkedItems returns all items linked to the given item.
func (r *ItemRelationRepo) GetLinkedItems(ctx context.Context, itemID uint64) ([]*entity.Item, error) {
	rows, err := r.queries.GetLinkedItems(ctx, int64(itemID))
	if errors.Is(err, sql.ErrNoRows) {
		return []*entity.Item{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ItemRelationRepo.GetLinkedItems: %w", err)
	}

	items := make([]*entity.Item, len(rows))
	for i, row := range rows {
		items[i] = mapLinkedItem(row)
	}
	return items, nil
}

// Delete removes a relation by its ID.
func (r *ItemRelationRepo) Delete(ctx context.Context, relationID uint64) error {
	err := r.queries.DeleteItemRelation(ctx, int64(relationID))
	if err != nil {
		return fmt.Errorf("ItemRelationRepo.Delete: %w", err)
	}
	return nil
}

// Exists checks if a relation exists between two items.
func (r *ItemRelationRepo) Exists(ctx context.Context, itemID1, itemID2 uint64) (bool, error) {
	// Enforce ordering for lookup
	if itemID1 > itemID2 {
		itemID1, itemID2 = itemID2, itemID1
	}

	exists, err := r.queries.RelationExists(ctx, db.RelationExistsParams{
		ItemID1: int64(itemID1),
		ItemID2: int64(itemID2),
	})
	if err != nil {
		return false, fmt.Errorf("ItemRelationRepo.Exists: %w", err)
	}
	return exists, nil
}

// mapLinkedItem converts a sqlc-generated GetLinkedItemsRow to a domain entity.
func mapLinkedItem(row db.GetLinkedItemsRow) *entity.Item {
	return &entity.Item{
		ID:              uint64(row.ID),
		Barcode:         row.Barcode,
		InventoryNumber: row.InventoryNumber,
		Name:            row.Name,

		CategoryID: uint64(row.CategoryID),
		Category: &entity.Category{
			ID:   uint64(row.CategoryID),
			Name: row.CategoryName,
		},

		RoomID: uint64(row.RoomID),
		Room: &entity.Room{
			ID:         uint64(row.RoomID),
			Name:       row.RoomName,
			BuildingID: uint64(row.BuildingID),
			Building: &entity.Building{
				ID:      uint64(row.BuildingID),
				Name:    row.BuildingName,
				Address: row.BuildingAddress,
			},
		},

		Description: row.Description,
		Status:      entity.ItemStatus(row.Status),
		TxHash:      row.TxHash,

		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,

		PendingDisposalAt: nullTimeToPtr(row.PendingDisposalAt),
		DisposedAt:        nullTimeToPtr(row.DisposedAt),

		CreatedBy: uint64(row.CreatedBy),
		Creator: &entity.User{
			ID:        uint64(row.CreatedBy),
			FullName:  row.CreatorFullName,
			Email:     row.CreatorEmail,
			Role:      entity.UserRole(row.CreatorRole),
			CreatedAt: row.CreatorCreatedAt,
			UpdatedAt: row.CreatorUpdatedAt,
		},

		LastEditedBy: uint64(row.LastEditedBy),
		LastEditor: &entity.User{
			ID:        uint64(row.LastEditedBy),
			FullName:  row.EditorFullName,
			Email:     row.EditorEmail,
			Role:      entity.UserRole(row.EditorRole),
			CreatedAt: row.EditorCreatedAt,
			UpdatedAt: row.EditorUpdatedAt,
		},
	}
}

// nullTimeToPtr converts sql.NullTime to *time.Time
func nullTimeToPtr(t sql.NullTime) *sql.NullTime {
	if !t.Valid {
		return nil
	}
	return &t
}
