// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// ItemRepo implements repository.ItemRepository using sqlc-generated queries.
type ItemRepo struct {
	queries *db.Queries
}

// NewItemRepo constructs a new ItemRepo backed by the given connection pool.
// pgxpool is adapted to database/sql via stdlib so sqlc-generated code can use it.
func NewItemRepo(pool *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Create inserts a new item row and populates the generated ID and timestamps.
func (r *ItemRepo) Create(ctx context.Context, item *entity.Item) error {
	row, err := r.queries.CreateItem(ctx, db.CreateItemParams{
		Barcode:     item.Barcode,
		Name:        item.Name,
		CategoryID:  int64(item.CategoryID),
		RoomID:      int64(item.RoomID),
		Description: item.Description,
		PhotoUrl:    item.PhotoURL,
		CreatedBy:   int64(item.CreatedBy),
	})
	if err != nil {
		return fmt.Errorf("ItemRepo.Create: %w", err)
	}

	item.ID = uint64(row.ID)
	item.CreatedAt = row.CreatedAt
	item.UpdatedAt = row.UpdatedAt
	return nil
}

// GetByID returns the item with the given ID or ErrNotFound.
func (r *ItemRepo) GetByID(ctx context.Context, id uint64) (*entity.Item, error) {
	row, err := r.queries.GetItemByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GetByID: %w", err)
	}
	return mapItemDetail(row), nil
}

// GetByBarcode returns the item matching the barcode or ErrNotFound.
func (r *ItemRepo) GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error) {
	row, err := r.queries.GetItemByBarcode(ctx, barcode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GetByBarcode: %w", err)
	}
	return mapItemDetail(row), nil
}

// List returns items matching the provided filter.
func (r *ItemRepo) List(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, error) {
	rows, err := r.queries.ListItems(ctx, db.ListItemsParams{
		CategoryID: nullInt64(filter.CategoryID),
		RoomID:     nullInt64(filter.RoomID),
		Status:     nullString((*string)(filter.Status)),
		DateFrom:   nullTime(filter.DateFrom),
		DateTo:     nullTime(filter.DateTo),
	})
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.List: %w", err)
	}

	items := make([]*entity.Item, len(rows))
	for i, row := range rows {
		items[i] = mapItemDetail(row)
	}
	return items, nil
}

// Update persists changes to Description, PhotoURL, and UpdatedAt.
func (r *ItemRepo) Update(ctx context.Context, item *entity.Item) error {
	err := r.queries.UpdateItem(ctx, db.UpdateItemParams{
		Description:  item.Description,
		PhotoUrl:     item.PhotoURL,
		LastEditedBy: int64(item.LastEditedBy),
		ID:           int64(item.ID),
	})
	if err != nil {
		return fmt.Errorf("ItemRepo.Update: %w", err)
	}
	return nil
}

// UpdateStatus persists a lifecycle status change (e.g. disposed).
func (r *ItemRepo) UpdateStatus(ctx context.Context, item *entity.Item) error {
	err := r.queries.UpdateItemStatus(ctx, db.UpdateItemStatusParams{
		Status:       string(item.Status),
		LastEditedBy: int64(item.LastEditedBy),
		ID:           int64(item.ID),
	})
	if err != nil {
		return fmt.Errorf("ItemRepo.UpdateStatus: %w", err)
	}
	return nil
}

// UpdateTxHash stores the blockchain transaction hash for an item.
func (r *ItemRepo) UpdateTxHash(ctx context.Context, id uint64, txHash string) error {
	err := r.queries.UpdateItemTxHash(ctx, db.UpdateItemTxHashParams{
		TxHash: txHash,
		ID:     int64(id),
	})
	if err != nil {
		return fmt.Errorf("ItemRepo.UpdateTxHash: %w", err)
	}
	return nil
}

// BarcodeExists reports whether the given barcode is already in use.
func (r *ItemRepo) BarcodeExists(ctx context.Context, barcode string) (bool, error) {
	exists, err := r.queries.BarcodeExists(ctx, barcode)
	if err != nil {
		return false, fmt.Errorf("ItemRepo.BarcodeExists: %w", err)
	}
	return exists, nil
}

// ── mapping ───────────────────────────────────────────────────────────────────

// mapItemDetail converts a sqlc-generated ItemDetail row to a domain entity.
// This is the single place where the DB representation maps to domain types.
func mapItemDetail(row db.ItemDetail) *entity.Item {
	return &entity.Item{
		ID:      uint64(row.ID),
		Barcode: row.Barcode,
		Name:    row.Name,

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
		PhotoURL:    row.PhotoUrl,
		Status:      entity.ItemStatus(row.Status),
		TxHash:      row.TxHash,

		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,

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

// ── sql.Null* helpers ─────────────────────────────────────────────────────────

// nullInt64 converts a nullable *uint64 filter to sql.NullInt64.
func nullInt64(v *uint64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

// nullString converts a nullable *string filter to sql.NullString.
func nullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

// nullTime converts a nullable *time.Time filter to sql.NullTime.
func nullTime(v *time.Time) sql.NullTime {
	if v == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *v, Valid: true}
}

