// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"net/http"
	"strconv"

	"github.com/Vallevas/Skopidom/internal/delivery/http/middleware"
	"github.com/Vallevas/Skopidom/internal/infrastructure/storage"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	"github.com/go-chi/chi/v5"
)

const maxUploadSize = 5 << 20 // 5 MB

// UploadHandler handles photo upload and deletion for inventory items.
type UploadHandler struct {
	storage storage.FileStorage
	itemUC  itemUC.UseCase
}

// NewUploadHandler constructs an UploadHandler.
func NewUploadHandler(s storage.FileStorage, uc itemUC.UseCase) *UploadHandler {
	return &UploadHandler{storage: s, itemUC: uc}
}

// UploadItemPhoto godoc
// POST /api/v1/items/{id}/photos
// Content-Type: multipart/form-data  field: "photo"
func (h *UploadHandler) UploadItemPhoto(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

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

	photo, err := h.itemUC.AddPhoto(r.Context(), itemID, photoURL, actorID)
	if err != nil {
		_ = h.storage.Delete(r.Context(), photoURL)
		handleError(w, err)
		return
	}

	respond(w, http.StatusCreated, photo)
}

// DeleteItemPhoto godoc
// DELETE /api/v1/items/{id}/photos/{photo_id}
func (h *UploadHandler) DeleteItemPhoto(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	photoIDStr := chi.URLParam(r, "photo_id")
	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	// Fetch the photo URL before deleting so we can remove the file.
	photos, err := h.itemUC.ListPhotos(r.Context(), itemID)
	if err != nil {
		handleError(w, err)
		return
	}

	var photoURL string
	for _, p := range photos {
		if p.ID == photoID {
			photoURL = p.URL
			break
		}
	}

	actorID := middleware.UserIDFromCtx(r.Context())
	if err := h.itemUC.DeletePhoto(r.Context(), itemID, photoID, actorID); err != nil {
		handleError(w, err)
		return
	}

	// Remove file from storage after successful DB deletion.
	if photoURL != "" {
		_ = h.storage.Delete(r.Context(), photoURL)
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemPhotos godoc
// GET /api/v1/items/{id}/photos
func (h *UploadHandler) ListItemPhotos(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	photos, err := h.itemUC.ListPhotos(r.Context(), itemID)
	if err != nil {
		handleError(w, err)
		return
	}

	respond(w, http.StatusOK, photos)
}
