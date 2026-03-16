// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// itemSelectBase is the single source of truth for the item SELECT query.
// All item queries (GetByID, GetByBarcode, List) use this base via queryItems.
const itemSelectBase = `
	SELECT
		i.id, i.barcode, i.name,
		i.category_id, c.name,
		i.room_id, rm.name,
		rm.building_id, b.name, b.address,
		i.description, i.photo_url, i.status, i.tx_hash,
		i.created_at, i.updated_at,
		i.created_by,
		uc.full_name, uc.email, uc.role, uc.created_at, uc.updated_at,
		i.last_edited_by,
		ue.full_name, ue.email, ue.role, ue.created_at, ue.updated_at
	FROM items i
	JOIN categories c  ON c.id  = i.category_id
	JOIN rooms      rm ON rm.id = i.room_id
	JOIN buildings  b  ON b.id  = rm.building_id
	JOIN users      uc ON uc.id = i.created_by
	JOIN users      ue ON ue.id = i.last_edited_by`

// ItemRepo implements repository.ItemRepository using PostgreSQL.
type ItemRepo struct {
	pool *pgxpool.Pool
}

// NewItemRepo constructs a new ItemRepo backed by the given connection pool.
func NewItemRepo(pool *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{pool: pool}
}

// Create inserts a new item row and populates the generated ID and timestamps.
func (r *ItemRepo) Create(ctx context.Context, item *entity.Item) error {
	query := `
		INSERT INTO items
			(barcode, name, category_id, room_id, description,
			 photo_url, status, created_by, last_edited_by)
		VALUES ($1, $2, $3, $4, $5, $6, 'active', $7, $7)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		item.Barcode, item.Name,
		item.CategoryID, item.RoomID,
		item.Description, item.PhotoURL,
		item.CreatedBy,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ItemRepo.Create: %w", err)
	}
	return nil
}

// GetByID returns the item with the given ID or ErrNotFound.
func (r *ItemRepo) GetByID(ctx context.Context, id uint64) (*entity.Item, error) {
	items, err := r.queryItems(ctx, "WHERE i.id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GetByID: %w", err)
	}
	if len(items) == 0 {
		return nil, logger.ErrNotFound
	}
	return items[0], nil
}

// GetByBarcode returns the item matching the barcode or ErrNotFound.
func (r *ItemRepo) GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error) {
	items, err := r.queryItems(ctx, "WHERE i.barcode = $1", barcode)
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.GetByBarcode: %w", err)
	}
	if len(items) == 0 {
		return nil, logger.ErrNotFound
	}
	return items[0], nil
}

// List returns items matching the provided filter.
func (r *ItemRepo) List(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, error) {
	where, args := buildItemWhere(filter)
	items, err := r.queryItems(ctx, where, args...)
	if err != nil {
		return nil, fmt.Errorf("ItemRepo.List: %w", err)
	}
	return items, nil
}

// Update persists changes to Description, PhotoURL, and UpdatedAt.
func (r *ItemRepo) Update(ctx context.Context, item *entity.Item) error {
	query := `
		UPDATE items
		SET description = $1, photo_url = $2, last_edited_by = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		item.Description, item.PhotoURL, item.LastEditedBy, item.ID,
	).Scan(&item.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("ItemRepo.Update: %w", err)
	}
	return nil
}

// UpdateStatus persists a lifecycle status change (e.g. disposed).
func (r *ItemRepo) UpdateStatus(ctx context.Context, item *entity.Item) error {
	query := `
		UPDATE items
		SET status = $1, last_edited_by = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query,
		item.Status, item.LastEditedBy, item.ID,
	).Scan(&item.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("ItemRepo.UpdateStatus: %w", err)
	}
	return nil
}

// UpdateTxHash stores the blockchain transaction hash for an item.
func (r *ItemRepo) UpdateTxHash(ctx context.Context, id uint64, txHash string) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE items SET tx_hash = $1 WHERE id = $2`, txHash, id,
	)
	if err != nil {
		return fmt.Errorf("ItemRepo.UpdateTxHash: %w", err)
	}
	if result.RowsAffected() == 0 {
		return logger.ErrNotFound
	}
	return nil
}

// BarcodeExists reports whether the given barcode is already in use.
func (r *ItemRepo) BarcodeExists(ctx context.Context, barcode string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM items WHERE barcode = $1)`, barcode,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ItemRepo.BarcodeExists: %w", err)
	}
	return exists, nil
}

// ── private helpers ───────────────────────────────────────────────────────────

// queryItems executes itemSelectBase with an optional WHERE clause and args,
// scans all rows, and returns the result slice.
// It is the single execution point for all item SELECT operations.
func (r *ItemRepo) queryItems(ctx context.Context, where string, args ...any) ([]*entity.Item, error) {
	query := itemSelectBase + "\n\t" + where + "\n\tORDER BY i.created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Item, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// buildItemWhere constructs a WHERE clause and argument slice from an ItemFilter.
// Returns an empty string and nil args if no filters are set.
func buildItemWhere(filter repository.ItemFilter) (string, []any) {
	args := make([]any, 0, 5)
	clause := "WHERE 1=1"
	idx := 1

	if filter.CategoryID != nil {
		clause += fmt.Sprintf(" AND i.category_id = $%d", idx)
		args = append(args, *filter.CategoryID)
		idx++
	}
	if filter.RoomID != nil {
		clause += fmt.Sprintf(" AND i.room_id = $%d", idx)
		args = append(args, *filter.RoomID)
		idx++
	}
	if filter.Status != nil {
		clause += fmt.Sprintf(" AND i.status = $%d", idx)
		args = append(args, *filter.Status)
		idx++
	}
	if filter.DateFrom != nil {
		clause += fmt.Sprintf(" AND i.created_at >= $%d", idx)
		args = append(args, *filter.DateFrom)
		idx++
	}
	if filter.DateTo != nil {
		clause += fmt.Sprintf(" AND i.created_at <= $%d", idx)
		args = append(args, *filter.DateTo)
		idx++
	}

	return clause, args
}

// rowScanner is satisfied by both pgx.Row and pgx.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanItem maps a single database row to a fully-populated Item value.
func scanItem(row rowScanner) (*entity.Item, error) {
	item := &entity.Item{
		Category:   &entity.Category{},
		Room:       &entity.Room{Building: &entity.Building{}},
		Creator:    &entity.User{},
		LastEditor: &entity.User{},
	}

	err := row.Scan(
		&item.ID, &item.Barcode, &item.Name,
		&item.CategoryID, &item.Category.Name,
		&item.RoomID, &item.Room.Name,
		&item.Room.BuildingID, &item.Room.Building.Name, &item.Room.Building.Address,
		&item.Description, &item.PhotoURL, &item.Status, &item.TxHash,
		&item.CreatedAt, &item.UpdatedAt,
		&item.CreatedBy,
		&item.Creator.FullName, &item.Creator.Email,
		&item.Creator.Role, &item.Creator.CreatedAt, &item.Creator.UpdatedAt,
		&item.LastEditedBy,
		&item.LastEditor.FullName, &item.LastEditor.Email,
		&item.LastEditor.Role, &item.LastEditor.CreatedAt, &item.LastEditor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.Category.ID = item.CategoryID
	item.Room.ID = item.RoomID
	item.Room.Building.ID = item.Room.BuildingID
	item.Creator.ID = item.CreatedBy
	item.LastEditor.ID = item.LastEditedBy

	return item, nil
}
