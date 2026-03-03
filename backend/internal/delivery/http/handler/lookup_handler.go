// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"context"
	"net/http"

	"fmt"
	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// ── Use case interfaces (thin, defined here to avoid extra packages) ──────────

// categoryUseCase is the minimal interface the handler depends on.
type categoryUseCase interface {
	Create(ctx context.Context, name string) (*entity.Category, error)
	List(ctx context.Context) ([]*entity.Category, error)
	Update(ctx context.Context, id uint64, name string) (*entity.Category, error)
	Delete(ctx context.Context, id uint64) error
}

// roomUseCase is the minimal interface the handler depends on.
type roomUseCase interface {
	Create(ctx context.Context, name string, buildingID uint64) (*entity.Room, error)
	List(ctx context.Context) ([]*entity.Room, error)
	ListByBuilding(ctx context.Context, buildingID uint64) ([]*entity.Room, error)
	Update(ctx context.Context, id uint64, name string, buildingID uint64) (*entity.Room, error)
	Delete(ctx context.Context, id uint64) error
}

// buildingUseCase is the minimal interface the handler depends on.
type buildingUseCase interface {
	Create(ctx context.Context, name, address string) (*entity.Building, error)
	List(ctx context.Context) ([]*entity.Building, error)
	Update(ctx context.Context, id uint64, name, address string) (*entity.Building, error)
	Delete(ctx context.Context, id uint64) error
}

// ── CategoryHandler ──────────────────────────────────────────────────────────

// CategoryHandler handles HTTP requests for category reference data.
type CategoryHandler struct {
	uc categoryUseCase
}

// NewCategoryHandler constructs a CategoryHandler.
func NewCategoryHandler(uc categoryUseCase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	cats, err := h.uc.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	cat, err := h.uc.Create(r.Context(), req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusCreated, cat)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	cat, err := h.uc.Update(r.Context(), id, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── BuildingHandler ──────────────────────────────────────────────────────────

// BuildingHandler handles HTTP requests for building reference data.
type BuildingHandler struct {
	uc buildingUseCase
}

// NewBuildingHandler constructs a BuildingHandler.
func NewBuildingHandler(uc buildingUseCase) *BuildingHandler {
	return &BuildingHandler{uc: uc}
}

func (h *BuildingHandler) List(w http.ResponseWriter, r *http.Request) {
	buildings, err := h.uc.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, buildings)
}

func (h *BuildingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	b, err := h.uc.Create(r.Context(), req.Name, req.Address)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusCreated, b)
}

func (h *BuildingHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	b, err := h.uc.Update(r.Context(), id, req.Name, req.Address)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, b)
}

func (h *BuildingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── RoomHandler ──────────────────────────────────────────────────────────────

// RoomHandler handles HTTP requests for room reference data.
type RoomHandler struct {
	uc roomUseCase
}

// NewRoomHandler constructs a RoomHandler.
func NewRoomHandler(uc roomUseCase) *RoomHandler {
	return &RoomHandler{uc: uc}
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	// Optional filter: ?building_id=1
	if raw := r.URL.Query().Get("building_id"); raw != "" {
		bid, err := parseUint(raw)
		if err != nil {
			handleError(w, wrapInvalidInput(err))
			return
		}
		rooms, err := h.uc.ListByBuilding(r.Context(), bid)
		if err != nil {
			handleError(w, err)
			return
		}
		respond(w, http.StatusOK, rooms)
		return
	}

	rooms, err := h.uc.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, rooms)
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		BuildingID uint64 `json:"building_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	room, err := h.uc.Create(r.Context(), req.Name, req.BuildingID)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusCreated, room)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	var req struct {
		Name       string `json:"name"`
		BuildingID uint64 `json:"building_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	room, err := h.uc.Update(r.Context(), id, req.Name, req.BuildingID)
	if err != nil {
		handleError(w, err)
		return
	}
	respond(w, http.StatusOK, room)
}

func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		handleError(w, wrapInvalidInput(err))
		return
	}
	if err := h.uc.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── shared helpers ────────────────────────────────────────────────────────────

func parseUint(s string) (uint64, error) {
	var val uint64
	_, err := fmt.Sscan(s, &val)
	return val, err
}

// lookupSimpleUseCase implements categoryUseCase / buildingUseCase via repos.
// This keeps use-case logic out of the handler package while avoiding an extra
// package for trivial CRUD operations.

// SimpleCategoryUC is a thin use case for categories backed directly by the repo.
type SimpleCategoryUC struct {
	repo repository.CategoryRepository
}

// NewSimpleCategoryUC constructs a SimpleCategoryUC.
func NewSimpleCategoryUC(repo repository.CategoryRepository) *SimpleCategoryUC {
	return &SimpleCategoryUC{repo: repo}
}

func (uc *SimpleCategoryUC) Create(ctx context.Context, name string) (*entity.Category, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", apperrors.ErrInvalidInput)
	}
	cat := &entity.Category{Name: name}
	if err := uc.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *SimpleCategoryUC) List(ctx context.Context) ([]*entity.Category, error) {
	return uc.repo.List(ctx)
}

func (uc *SimpleCategoryUC) Update(ctx context.Context, id uint64, name string) (*entity.Category, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", apperrors.ErrInvalidInput)
	}
	cat := &entity.Category{ID: id, Name: name}
	if err := uc.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (uc *SimpleCategoryUC) Delete(ctx context.Context, id uint64) error {
	return uc.repo.Delete(ctx, id)
}

// SimpleBuildingUC is a thin use case for buildings backed directly by the repo.
type SimpleBuildingUC struct {
	repo repository.BuildingRepository
}

// NewSimpleBuildingUC constructs a SimpleBuildingUC.
func NewSimpleBuildingUC(repo repository.BuildingRepository) *SimpleBuildingUC {
	return &SimpleBuildingUC{repo: repo}
}

func (uc *SimpleBuildingUC) Create(ctx context.Context, name, address string) (*entity.Building, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", apperrors.ErrInvalidInput)
	}
	b := &entity.Building{Name: name, Address: address}
	if err := uc.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (uc *SimpleBuildingUC) List(ctx context.Context) ([]*entity.Building, error) {
	return uc.repo.List(ctx)
}

func (uc *SimpleBuildingUC) Update(ctx context.Context, id uint64, name, address string) (*entity.Building, error) {
	b := &entity.Building{ID: id, Name: name, Address: address}
	if err := uc.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (uc *SimpleBuildingUC) Delete(ctx context.Context, id uint64) error {
	return uc.repo.Delete(ctx, id)
}

// SimpleRoomUC is a thin use case for rooms backed directly by the repo.
type SimpleRoomUC struct {
	repo      repository.RoomRepository
	buildings repository.BuildingRepository
}

// NewSimpleRoomUC constructs a SimpleRoomUC.
func NewSimpleRoomUC(
	repo repository.RoomRepository,
	buildings repository.BuildingRepository,
) *SimpleRoomUC {
	return &SimpleRoomUC{repo: repo, buildings: buildings}
}

func (uc *SimpleRoomUC) Create(ctx context.Context, name string, buildingID uint64) (*entity.Room, error) {
	if name == "" || buildingID == 0 {
		return nil, fmt.Errorf("name and building_id are required: %w", apperrors.ErrInvalidInput)
	}
	// Verify building exists.
	if _, err := uc.buildings.GetByID(ctx, buildingID); err != nil {
		return nil, fmt.Errorf("building: %w", err)
	}
	room := &entity.Room{Name: name, BuildingID: buildingID}
	if err := uc.repo.Create(ctx, room); err != nil {
		return nil, err
	}
	// Re-fetch to get the joined building data.
	return uc.repo.GetByID(ctx, room.ID)
}

func (uc *SimpleRoomUC) List(ctx context.Context) ([]*entity.Room, error) {
	return uc.repo.List(ctx)
}

func (uc *SimpleRoomUC) ListByBuilding(ctx context.Context, buildingID uint64) ([]*entity.Room, error) {
	return uc.repo.ListByBuilding(ctx, buildingID)
}

func (uc *SimpleRoomUC) Update(ctx context.Context, id uint64, name string, buildingID uint64) (*entity.Room, error) {
	room := &entity.Room{ID: id, Name: name, BuildingID: buildingID}
	if err := uc.repo.Update(ctx, room); err != nil {
		return nil, err
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *SimpleRoomUC) Delete(ctx context.Context, id uint64) error {
	return uc.repo.Delete(ctx, id)
}
