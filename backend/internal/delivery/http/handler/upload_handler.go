// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"encoding/base64"
	"io"
	"net/http"
	"strconv"

	"github.com/Vallevas/Skopidom/internal/delivery/http/middleware"
	"github.com/Vallevas/Skopidom/internal/domain/entity"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	"github.com/go-chi/chi/v5"
)

const maxUploadSize = 5 << 20 // 5 MB

// UploadHandler handles photo upload and deletion for inventory items.
type UploadHandler struct {
	itemUC itemUC.UseCase
}

// NewUploadHandler constructs an UploadHandler.
func NewUploadHandler(s interface{}, uc itemUC.UseCase) *UploadHandler {
	return &UploadHandler{itemUC: uc}
}

// detectContentType reads up to 512 bytes from the file to detect MIME type.
func detectContentType(file io.Reader) (string, []byte, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", nil, err
	}
	contentType := http.DetectContentType(buf[:n])
	// Read the rest of the file
	rest, err := io.ReadAll(file)
	if err != nil {
		return "", nil, err
	}
	data := append(buf[:n], rest...)
	return contentType, data, nil
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

	// Detect MIME type and read file content
	mimeType, fileData, err := detectContentType(file)
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	// Encode to Base64
	base64Data := base64.StdEncoding.EncodeToString(fileData)

	actorID := middleware.UserIDFromCtx(r.Context())

	photo := &entity.ItemPhoto{
		ItemID:     itemID,
		Base64Data: base64Data,
		MimeType:   mimeType,
	}

	photo, err = h.itemUC.AddPhotoWithEntity(r.Context(), itemID, photo, actorID)
	if err != nil {
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

	actorID := middleware.UserIDFromCtx(r.Context())
	if err := h.itemUC.DeletePhoto(r.Context(), itemID, photoID, actorID); err != nil {
		handleError(w, err)
		return
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

// UploadDisposalDocument godoc
// POST /api/v1/items/{id}/disposal-documents
// Content-Type: multipart/form-data  field: "document"
func (h *UploadHandler) UploadDisposalDocument(w http.ResponseWriter, r *http.Request) {
	itemID, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(config.MaxDisposalDocumentSize))
	if err := r.ParseMultipartForm(int64(config.MaxDisposalDocumentSize)); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	defer file.Close()

	docURL, err := h.storage.SaveDocument(r.Context(), file, header)
	if err != nil {
		handleError(w, err)
		return
	}

	actorID := middleware.UserIDFromCtx(r.Context())

	doc, err := h.itemUC.UploadDisposalDocument(r.Context(), itemID, header.Filename, docURL, actorID)
	if err != nil {
		_ = h.storage.Delete(r.Context(), docURL)
		handleError(w, err)
		return
	}

	respond(w, http.StatusCreated, doc)
}
