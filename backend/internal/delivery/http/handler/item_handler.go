// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Vallevas/Skopidom/internal/delivery/http/middleware"
	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	"github.com/go-chi/chi/v5"
)

// ItemHandler handles HTTP requests for inventory item resources.
type ItemHandler struct {
	uc itemUC.UseCase
}

// NewItemHandler constructs an ItemHandler with the given use case.
func NewItemHandler(uc itemUC.UseCase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

// ── Request / Response types ──────────────────────────────────────────────────

type createItemRequest struct {
	Barcode     string `json:"barcode"`
	Name        string `json:"name"`
	CategoryID  uint64 `json:"category_id"`
	RoomID      uint64 `json:"room_id"`
	Description string `json:"description"`
}

type updateItemRequest struct {
	Description string `json:"description"`
}

type moveToRoomRequest struct {
	RoomID uint64 `json:"room_id"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// Create godoc
// POST /api/v1/items
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}

	actorID := middleware.UserIDFromCtx(r.Context())
	item, err := h.uc.Create(r.Context(), itemUC.CreateInput{
		Barcode:     req.Barcode,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		RoomID:      req.RoomID,
		Description: req.Description,
		ActorID:     actorID,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusCreated, item)
}

// GetByID godoc
// GET /api/v1/items/{id}
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	item, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, item)
}

// GetByBarcode godoc
// GET /api/v1/items/barcode/{barcode}
func (h *ItemHandler) GetByBarcode(w http.ResponseWriter, r *http.Request) {
	barcode := chi.URLParam(r, "barcode")
	item, err := h.uc.GetByBarcode(r.Context(), barcode)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, item)
}

// List godoc
// GET /api/v1/items?category_id=&room_id=&status=&date_from=&date_to=
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := buildItemFilter(r)
	items, err := h.uc.List(r.Context(), filter)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, items)
}

// Update godoc
// PATCH /api/v1/items/{id}
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	var req updateItemRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	actorID := middleware.UserIDFromCtx(r.Context())
	item, err := h.uc.Update(r.Context(), itemUC.UpdateInput{
		ItemID:      id,
		Description: req.Description,
		ActorID:     actorID,
	})
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, item)
}

// MoveToRoom godoc
// PATCH /api/v1/items/{id}/room
func (h *ItemHandler) MoveToRoom(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	var req moveToRoomRequest
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	actorID := middleware.UserIDFromCtx(r.Context())
	item, err := h.uc.MoveToRoom(r.Context(), id, req.RoomID, actorID)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, item)
}

// ToggleRepair godoc
// PATCH /api/v1/items/{id}/repair
// Switches active ↔ in_repair. No request body needed.
func (h *ItemHandler) ToggleRepair(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	actorID := middleware.UserIDFromCtx(r.Context())
	item, err := h.uc.ToggleRepair(r.Context(), id, actorID)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, item)
}

// GetAuditLog godoc
// GET /api/v1/items/{id}/audit
func (h *ItemHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	events, err := h.uc.ListAuditEvents(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, events)
}

// Dispose godoc
// DELETE /api/v1/items/{id}  (admin only)
func (h *ItemHandler) Dispose(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	actorID := middleware.UserIDFromCtx(r.Context())
	if err := h.uc.Dispose(r.Context(), id, actorID); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseIDParam(r *http.Request, name string) (uint64, error) {
	return strconv.ParseUint(chi.URLParam(r, name), 10, 64)
}

func buildItemFilter(r *http.Request) repository.ItemFilter {
	q := r.URL.Query()
	filter := repository.ItemFilter{}

	if raw := q.Get("category_id"); raw != "" {
		if val, err := strconv.ParseUint(raw, 10, 64); err == nil {
			filter.CategoryID = &val
		}
	}
	if raw := q.Get("room_id"); raw != "" {
		if val, err := strconv.ParseUint(raw, 10, 64); err == nil {
			filter.RoomID = &val
		}
	}
	if raw := q.Get("status"); raw != "" {
		status := entity.ItemStatus(raw)
		filter.Status = &status
	}
	if raw := q.Get("date_from"); raw != "" {
		if val, err := time.Parse(time.RFC3339, raw); err == nil {
			filter.DateFrom = &val
		}
	}
	if raw := q.Get("date_to"); raw != "" {
		if val, err := time.Parse(time.RFC3339, raw); err == nil {
			filter.DateTo = &val
		}
	}
	return filter
}
