// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// PhotoRepo implements repository.PhotoRepository using sqlc-generated queries.
type PhotoRepo struct {
	queries *db.Queries
}

// NewPhotoRepo constructs a PhotoRepo backed by the given connection pool.
func NewPhotoRepo(pool *pgxpool.Pool) *PhotoRepo {
	return &PhotoRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Add inserts a new photo record and populates the generated ID and timestamp.
func (r *PhotoRepo) Add(ctx context.Context, photo *entity.ItemPhoto) error {
	row, err := r.queries.AddItemPhoto(ctx, db.AddItemPhotoParams{
		ItemID: int64(photo.ItemID),
		Url:    photo.URL,
	})
	if err != nil {
		return fmt.Errorf("PhotoRepo.Add: %w", err)
	}
	photo.ID = uint64(row.ID)
	photo.CreatedAt = row.CreatedAt
	return nil
}

// ListByItem returns all photos for the given item ordered by upload time.
func (r *PhotoRepo) ListByItem(ctx context.Context, itemID uint64) ([]*entity.ItemPhoto, error) {
	rows, err := r.queries.ListItemPhotos(ctx, int64(itemID))
	if err != nil {
		return nil, fmt.Errorf("PhotoRepo.ListByItem: %w", err)
	}
	photos := make([]*entity.ItemPhoto, len(rows))
	for i, row := range rows {
		photos[i] = &entity.ItemPhoto{
			ID:        uint64(row.ID),
			ItemID:    uint64(row.ItemID),
			URL:       row.Url,
			CreatedAt: row.CreatedAt,
		}
	}
	return photos, nil
}

// Delete removes a single photo verifying it belongs to the given item.
func (r *PhotoRepo) Delete(ctx context.Context, photoID uint64, itemID uint64) error {
	err := r.queries.DeleteItemPhoto(ctx, db.DeleteItemPhotoParams{
		ID:     int64(photoID),
		ItemID: int64(itemID),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("PhotoRepo.Delete: %w", err)
	}
	return nil
}
