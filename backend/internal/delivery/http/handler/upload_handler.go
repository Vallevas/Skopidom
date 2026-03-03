// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"net/http"

	"github.com/Vallevas/Skopidom/internal/delivery/http/middleware"
	"github.com/Vallevas/Skopidom/internal/infrastructure/storage"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
)

const (
	// maxUploadSize limits photo uploads to 5 MB.
	maxUploadSize = 5 << 20
)

// UploadHandler handles photo upload and attachment to inventory items.
type UploadHandler struct {
	storage storage.FileStorage
	itemUC  itemUC.UseCase
}

// NewUploadHandler constructs an UploadHandler.
func NewUploadHandler(s storage.FileStorage, uc itemUC.UseCase) *UploadHandler {
	return &UploadHandler{storage: s, itemUC: uc}
}

// UploadItemPhoto godoc
// POST /api/v1/items/{id}/photo
// Content-Type: multipart/form-data  field: "photo"
func (h *UploadHandler) UploadItemPhoto(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	// Enforce size limit before parsing the multipart body.
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	defer file.Close()

	photoURL, err := h.storage.Save(r.Context(), file, header)
	if err != nil {
		handleError(w, err)
		return
	}

	actorID := middleware.UserIDFromCtx(r.Context())

	// Fetch current description so the Update call does not overwrite it.
	current, err := h.itemUC.GetByID(r.Context(), itemID)
	if err != nil {
		// Remove the saved file to avoid orphaned uploads.
		_ = h.storage.Delete(r.Context(), photoURL)
		handleError(w, err)
		return
	}

	// If the item already has a photo, delete the old file.
	if current.PhotoURL != "" {
		_ = h.storage.Delete(r.Context(), current.PhotoURL)
	}

	updated, err := h.itemUC.Update(r.Context(), itemUC.UpdateInput{
		ItemID:      itemID,
		Description: current.Description,
		PhotoURL:    photoURL,
		ActorID:     actorID,
	})
	if err != nil {
		_ = h.storage.Delete(r.Context(), photoURL)
		handleError(w, err)
		return
	}

	respond(w, http.StatusOK, updated)
}
