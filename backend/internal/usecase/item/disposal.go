// Package item contains the business logic for inventory item disposal.
package item

import (
	"context"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/pkg/config"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// InitiateDisposal transitions an item to pending_disposal status.
// Only active and in_repair items can be disposed.
func (uc *itemUseCase) InitiateDisposal(ctx context.Context, itemID uint64, actorID uint64) (*entity.Item, error) {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Only active and in_repair items can be disposed.
	if item.Status != entity.StatusActive && item.Status != entity.StatusInRepair {
		return nil, logger.NewBusinessError("item must be active or in_repair to initiate disposal")
	}

	// Transition to pending_disposal.
	item.MarkPendingDisposal(actorID)

	if err := uc.items.UpdateStatus(ctx, item); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, item, entity.ActionPendingDisposal, actorID)

	return uc.items.GetByID(ctx, itemID)
}

// UploadDisposalDocument uploads a disposal document for an item in pending_disposal status.
func (uc *itemUseCase) UploadDisposalDocument(ctx context.Context, itemID uint64, filename string, url string, actorID uint64) (*entity.DisposalDocument, error) {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Only items in pending_disposal can have documents uploaded.
	if item.Status != entity.StatusPendingDisposal {
		return nil, logger.NewBusinessError("item must be in pending_disposal status to upload documents")
	}

	// Check document count limit.
	count, err := uc.disposalDocs.CountByItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if count >= int64(config.MaxDisposalDocumentsPerItem) {
		return nil, logger.NewBusinessError(fmt.Sprintf("maximum %d disposal documents allowed per item", config.MaxDisposalDocumentsPerItem))
	}

	doc := &entity.DisposalDocument{
		ItemID:     itemID,
		Filename:   filename,
		URL:        url,
		UploadedBy: actorID,
	}

	if err := uc.disposalDocs.Create(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ListDisposalDocuments returns all disposal documents for an item.
func (uc *itemUseCase) ListDisposalDocuments(ctx context.Context, itemID uint64) ([]*entity.DisposalDocument, error) {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Only items in pending_disposal or disposed can have documents listed.
	if item.Status != entity.StatusPendingDisposal && item.Status != entity.StatusDisposed {
		return nil, logger.NewBusinessError("item must be in pending_disposal or disposed status")
	}

	return uc.disposalDocs.ListByItem(ctx, itemID)
}

// DeleteDisposalDocument removes a disposal document (only allowed in pending_disposal status).
func (uc *itemUseCase) DeleteDisposalDocument(ctx context.Context, itemID uint64, docID uint64, actorID uint64) error {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return err
	}

	// Only items in pending_disposal can have documents deleted.
	if item.Status != entity.StatusPendingDisposal {
		return logger.NewBusinessError("can only delete documents while item is in pending_disposal status")
	}

	doc, err := uc.disposalDocs.GetByID(ctx, docID)
	if err != nil {
		return err
	}

	if doc.ItemID != itemID {
		return logger.NewBusinessError("document does not belong to this item")
	}

	return uc.disposalDocs.Delete(ctx, docID, itemID)
}

// FinalizeDisposal completes the disposal process, transitioning item to disposed status.
// Requires at least one disposal document to be uploaded.
func (uc *itemUseCase) FinalizeDisposal(ctx context.Context, itemID uint64, actorID uint64) (*entity.Item, error) {
	item, err := uc.items.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Only items in pending_disposal can be finalized.
	if item.Status != entity.StatusPendingDisposal {
		return nil, logger.NewBusinessError("item must be in pending_disposal status to finalize disposal")
	}

	// Check that at least one document is uploaded.
	count, err := uc.disposalDocs.CountByItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, logger.NewBusinessError("at least one disposal document must be uploaded before finalizing")
	}

	// Transition to disposed.
	item.FinalizeDisposal(actorID)

	if err := uc.items.UpdateStatus(ctx, item); err != nil {
		return nil, err
	}

	uc.logEvent(ctx, item, entity.ActionDisposalFinalized, actorID)

	return uc.items.GetByID(ctx, itemID)
}
