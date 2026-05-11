// Package item contains the business logic for inventory item management.
package item

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
)

// UseCase defines all business operations on inventory items.
type UseCase interface {
	Create(ctx context.Context, input CreateInput) (*entity.Item, error)
	Update(ctx context.Context, input UpdateInput) (*entity.Item, error)
	GetByID(ctx context.Context, id uint64) (*entity.Item, error)
	GetByBarcode(ctx context.Context, barcode string) (*entity.Item, error)
	List(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, error)
	ListAuditEvents(ctx context.Context, itemID uint64) ([]*entity.AuditEvent, error)
	MoveToRoom(ctx context.Context, itemID uint64, roomID uint64, actorID uint64) (*entity.Item, error)
	AddPhoto(ctx context.Context, itemID uint64, url string, actorID uint64) (*entity.ItemPhoto, error)
	ListPhotos(ctx context.Context, itemID uint64) ([]*entity.ItemPhoto, error)
	DeletePhoto(ctx context.Context, itemID uint64, photoID uint64, actorID uint64) error

	// ToggleRepair switches an active item to in_repair and vice versa.
	// Returns the updated item and the audit action that was recorded.
	ToggleRepair(ctx context.Context, itemID uint64, actorID uint64) (*entity.Item, error)

	// Disposal workflow methods.
	InitiateDisposal(ctx context.Context, itemID uint64, actorID uint64) (*entity.Item, error)
	UploadDisposalDocument(ctx context.Context, itemID uint64, filename string, url string, actorID uint64) (*entity.DisposalDocument, error)
	ListDisposalDocuments(ctx context.Context, itemID uint64) ([]*entity.DisposalDocument, error)
	DeleteDisposalDocument(ctx context.Context, itemID uint64, docID uint64, actorID uint64) error
	FinalizeDisposal(ctx context.Context, itemID uint64, actorID uint64) (*entity.Item, error)

	// Photo management with entity
	AddPhotoWithEntity(ctx context.Context, itemID uint64, photo *entity.ItemPhoto, actorID uint64) (*entity.ItemPhoto, error)
}

// CreateInput holds the data required to register a new inventory item.
type CreateInput struct {
	Barcode         string
	InventoryNumber string
	Name            string
	CategoryID      uint64
	RoomID          uint64
	Description     string
	ActorID         uint64
}

// UpdateInput holds the data allowed to be changed after creation.
type UpdateInput struct {
	ItemID      uint64
	Description string
	ActorID     uint64
}

// itemUseCase is the concrete implementation of UseCase.
type itemUseCase struct {
	items        repository.ItemRepository
	categories   repository.CategoryRepository
	rooms        repository.RoomRepository
	photos       repository.PhotoRepository
	disposalDocs repository.DisposalDocumentRepository
	audit        repository.AuditLogger
}

// New constructs an itemUseCase with all required repository dependencies.
func New(
	items repository.ItemRepository,
	categories repository.CategoryRepository,
	rooms repository.RoomRepository,
	photos repository.PhotoRepository,
	disposalDocs repository.DisposalDocumentRepository,
	audit repository.AuditLogger,
) UseCase {
	return &itemUseCase{
		items:        items,
		categories:   categories,
		rooms:        rooms,
		photos:       photos,
		disposalDocs: disposalDocs,
		audit:        audit,
	}
}

// ListAuditEvents returns the full audit history for the given item.
func (uc *itemUseCase) ListAuditEvents(ctx context.Context, itemID uint64) ([]*entity.AuditEvent, error) {
	return uc.audit.ListByItem(ctx, itemID)
}

// ListPhotos returns all photos for an item.
func (uc *itemUseCase) ListPhotos(ctx context.Context, itemID uint64) ([]*entity.ItemPhoto, error) {
	return uc.photos.ListByItem(ctx, itemID)
}

// logEvent builds and persists an AuditEvent — never fails the caller.
func (uc *itemUseCase) logEvent(
	ctx context.Context,
	item *entity.Item,
	action entity.AuditAction,
	actorID uint64,
) {
	payload, err := json.Marshal(item)
	if err != nil {
		slog.Error("audit: failed to marshal item snapshot",
			"item_id", item.ID, "err", err)
		return
	}
	_ = uc.audit.Log(ctx, &entity.AuditEvent{
		ItemID:  item.ID,
		ActorID: actorID,
		Action:  action,
		Payload: string(payload),
	})
}
