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

// DisposalDocumentRepo implements repository.DisposalDocumentRepository using sqlc-generated queries.
type DisposalDocumentRepo struct {
	queries *db.Queries
}

// NewDisposalDocumentRepo constructs a DisposalDocumentRepo backed by the given connection pool.
func NewDisposalDocumentRepo(pool *pgxpool.Pool) *DisposalDocumentRepo {
	return &DisposalDocumentRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Create inserts a new disposal document record and populates the generated ID and timestamp.
func (r *DisposalDocumentRepo) Create(ctx context.Context, doc *entity.DisposalDocument) error {
	row, err := r.queries.CreateDisposalDocument(ctx, db.CreateDisposalDocumentParams{
		ItemID:     int64(doc.ItemID),
		Filename:   doc.Filename,
		Url:        doc.URL,
		UploadedBy: int64(doc.UploadedBy),
	})
	if err != nil {
		return fmt.Errorf("DisposalDocumentRepo.Create: %w", err)
	}
	doc.ID = uint64(row.ID)
	doc.UploadedAt = row.UploadedAt
	return nil
}

// GetByID retrieves a single disposal document by its ID.
func (r *DisposalDocumentRepo) GetByID(ctx context.Context, id uint64) (*entity.DisposalDocument, error) {
	row, err := r.queries.GetDisposalDocumentByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("DisposalDocumentRepo.GetByID: %w", err)
	}
	return &entity.DisposalDocument{
		ID:         uint64(row.ID),
		ItemID:     uint64(row.ItemID),
		Filename:   row.Filename,
		URL:        row.Url,
		UploadedAt: row.UploadedAt,
		UploadedBy: uint64(row.UploadedBy),
	}, nil
}

// ListByItem returns all disposal documents for the given item.
func (r *DisposalDocumentRepo) ListByItem(ctx context.Context, itemID uint64) ([]*entity.DisposalDocument, error) {
	rows, err := r.queries.ListDisposalDocumentsByItemID(ctx, int64(itemID))
	if err != nil {
		return nil, fmt.Errorf("DisposalDocumentRepo.ListByItem: %w", err)
	}
	docs := make([]*entity.DisposalDocument, len(rows))
	for i, row := range rows {
		docs[i] = &entity.DisposalDocument{
			ID:         uint64(row.ID),
			ItemID:     uint64(row.ItemID),
			Filename:   row.Filename,
			URL:        row.Url,
			UploadedAt: row.UploadedAt,
			UploadedBy: uint64(row.UploadedBy),
		}
	}
	return docs, nil
}

// CountByItem returns the number of disposal documents for the given item.
func (r *DisposalDocumentRepo) CountByItem(ctx context.Context, itemID uint64) (int64, error) {
	count, err := r.queries.CountDisposalDocumentsByItemID(ctx, int64(itemID))
	if err != nil {
		return 0, fmt.Errorf("DisposalDocumentRepo.CountByItem: %w", err)
	}
	return count, nil
}

// Delete removes a disposal document by its ID, verifying it belongs to the given item.
func (r *DisposalDocumentRepo) Delete(ctx context.Context, docID uint64, itemID uint64) error {
	err := r.queries.DeleteDisposalDocument(ctx, db.DeleteDisposalDocumentParams{
		ID:     int64(docID),
		ItemID: int64(itemID),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return logger.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("DisposalDocumentRepo.Delete: %w", err)
	}
	return nil
}
