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

// ── CategoryRepo ──────────────────────────────────────────────────────────────

// CategoryRepo implements repository.CategoryRepository using sqlc-generated queries.
type CategoryRepo struct {
	queries *db.Queries
}

// NewCategoryRepo constructs a CategoryRepo backed by the given pool.
func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

func (r *CategoryRepo) Create(ctx context.Context, cat *entity.Category) error {
	id, err := r.queries.CreateCategory(ctx, cat.Name)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Create: %w", err)
	}
	cat.ID = uint64(id)
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uint64) (*entity.Category, error) {
	row, err := r.queries.GetCategoryByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.GetByID: %w", err)
	}
	return &entity.Category{ID: uint64(row.ID), Name: row.Name}, nil
}

func (r *CategoryRepo) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	row, err := r.queries.GetCategoryByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.GetByName: %w", err)
	}
	return &entity.Category{ID: uint64(row.ID), Name: row.Name}, nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]*entity.Category, error) {
	rows, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.List: %w", err)
	}
	cats := make([]*entity.Category, len(rows))
	for i, row := range rows {
		cats[i] = &entity.Category{ID: uint64(row.ID), Name: row.Name}
	}
	return cats, nil
}

func (r *CategoryRepo) Update(ctx context.Context, cat *entity.Category) error {
	err := r.queries.UpdateCategory(ctx, db.UpdateCategoryParams{
		Name: cat.Name,
		ID:   int64(cat.ID),
	})
	if err != nil {
		return fmt.Errorf("CategoryRepo.Update: %w", err)
	}
	return nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id uint64) error {
	// Check if category has items
	count, err := r.queries.CountItemsByCategory(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("CategoryRepo.Delete: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category: %d items are using this category: %w", count, logger.ErrConflict)
	}

	if err := r.queries.DeleteCategory(ctx, int64(id)); err != nil {
		return fmt.Errorf("CategoryRepo.Delete: %w", err)
	}
	return nil
}

// ── BuildingRepo ──────────────────────────────────────────────────────────────

// BuildingRepo implements repository.BuildingRepository using sqlc-generated queries.
type BuildingRepo struct {
	queries *db.Queries
}

// NewBuildingRepo constructs a BuildingRepo backed by the given pool.
func NewBuildingRepo(pool *pgxpool.Pool) *BuildingRepo {
	return &BuildingRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

func (r *BuildingRepo) Create(ctx context.Context, b *entity.Building) error {
	id, err := r.queries.CreateBuilding(ctx, db.CreateBuildingParams{
		Name:    b.Name,
		Address: b.Address,
	})
	if err != nil {
		return fmt.Errorf("BuildingRepo.Create: %w", err)
	}
	b.ID = uint64(id)
	return nil
}

func (r *BuildingRepo) GetByID(ctx context.Context, id uint64) (*entity.Building, error) {
	row, err := r.queries.GetBuildingByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("BuildingRepo.GetByID: %w", err)
	}
	return &entity.Building{ID: uint64(row.ID), Name: row.Name, Address: row.Address}, nil
}

func (r *BuildingRepo) GetByName(ctx context.Context, name string) (*entity.Building, error) {
	row, err := r.queries.GetBuildingByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("BuildingRepo.GetByName: %w", err)
	}
	return &entity.Building{ID: uint64(row.ID), Name: row.Name, Address: row.Address}, nil
}

func (r *BuildingRepo) List(ctx context.Context) ([]*entity.Building, error) {
	rows, err := r.queries.ListBuildings(ctx)
	if err != nil {
		return nil, fmt.Errorf("BuildingRepo.List: %w", err)
	}
	buildings := make([]*entity.Building, len(rows))
	for i, row := range rows {
		buildings[i] = &entity.Building{
			ID:      uint64(row.ID),
			Name:    row.Name,
			Address: row.Address,
		}
	}
	return buildings, nil
}

func (r *BuildingRepo) Update(ctx context.Context, b *entity.Building) error {
	err := r.queries.UpdateBuilding(ctx, db.UpdateBuildingParams{
		Name:    b.Name,
		Address: b.Address,
		ID:      int64(b.ID),
	})
	if err != nil {
		return fmt.Errorf("BuildingRepo.Update: %w", err)
	}
	return nil
}

func (r *BuildingRepo) Delete(ctx context.Context, id uint64) error {
	// Check if building has rooms
	count, err := r.queries.CountRoomsByBuilding(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("BuildingRepo.Delete: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete building: %d rooms are in this building: %w", count, logger.ErrConflict)
	}

	if err := r.queries.DeleteBuilding(ctx, int64(id)); err != nil {
		return fmt.Errorf("BuildingRepo.Delete: %w", err)
	}
	return nil
}

// ── RoomRepo ──────────────────────────────────────────────────────────────────

// RoomRepo implements repository.RoomRepository using sqlc-generated queries.
type RoomRepo struct {
	queries *db.Queries
}

// NewRoomRepo constructs a RoomRepo backed by the given pool.
func NewRoomRepo(pool *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

func (r *RoomRepo) Create(ctx context.Context, room *entity.Room) error {
	id, err := r.queries.CreateRoom(ctx, db.CreateRoomParams{
		Name:       room.Name,
		BuildingID: int64(room.BuildingID),
	})
	if err != nil {
		return fmt.Errorf("RoomRepo.Create: %w", err)
	}
	room.ID = uint64(id)
	return nil
}

func (r *RoomRepo) GetByID(ctx context.Context, id uint64) (*entity.Room, error) {
	row, err := r.queries.GetRoomByID(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.GetByID: %w", err)
	}
	return mapRoomFromGetByID(row), nil
}

func (r *RoomRepo) GetByNameAndBuilding(ctx context.Context, name string, buildingID uint64) (*entity.Room, error) {
	row, err := r.queries.GetRoomByNameAndBuilding(ctx, db.GetRoomByNameAndBuildingParams{
		Name:       name,
		BuildingID: int64(buildingID),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logger.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.GetByNameAndBuilding: %w", err)
	}
	return mapRoomFromGetByNameAndBuilding(row), nil
}

func (r *RoomRepo) List(ctx context.Context) ([]*entity.Room, error) {
	rows, err := r.queries.ListRooms(ctx)
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.List: %w", err)
	}
	rooms := make([]*entity.Room, len(rows))
	for i, row := range rows {
		rooms[i] = mapRoomFromList(row)
	}
	return rooms, nil
}

func (r *RoomRepo) ListByBuilding(ctx context.Context, buildingID uint64) ([]*entity.Room, error) {
	rows, err := r.queries.ListRoomsByBuilding(ctx, int64(buildingID))
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.ListByBuilding: %w", err)
	}
	rooms := make([]*entity.Room, len(rows))
	for i, row := range rows {
		rooms[i] = mapRoomFromListByBuilding(row)
	}
	return rooms, nil
}

func (r *RoomRepo) Update(ctx context.Context, room *entity.Room) error {
	err := r.queries.UpdateRoom(ctx, db.UpdateRoomParams{
		Name:       room.Name,
		BuildingID: int64(room.BuildingID),
		ID:         int64(room.ID),
	})
	if err != nil {
		return fmt.Errorf("RoomRepo.Update: %w", err)
	}
	return nil
}

func (r *RoomRepo) Delete(ctx context.Context, id uint64) error {
	// Check if room has items
	count, err := r.queries.CountItemsByRoom(ctx, int64(id))
	if err != nil {
		return fmt.Errorf("RoomRepo.Delete: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete room: %d items are in this room: %w", count, logger.ErrConflict)
	}

	if err := r.queries.DeleteRoom(ctx, int64(id)); err != nil {
		return fmt.Errorf("RoomRepo.Delete: %w", err)
	}
	return nil
}

// ── mapping ───────────────────────────────────────────────────────────────────
// sqlc generates a distinct row type per query even when columns are identical.
// Three small mappers avoid duplicating the field mapping logic.

func mapRoomFromGetByID(row db.GetRoomByIDRow) *entity.Room {
	return &entity.Room{
		ID:         uint64(row.ID),
		Name:       row.Name,
		BuildingID: uint64(row.BuildingID),
		Building: &entity.Building{
			ID:      uint64(row.BuildingID),
			Name:    row.BuildingName,
			Address: row.BuildingAddress,
		},
	}
}

func mapRoomFromList(row db.ListRoomsRow) *entity.Room {
	return &entity.Room{
		ID:         uint64(row.ID),
		Name:       row.Name,
		BuildingID: uint64(row.BuildingID),
		Building: &entity.Building{
			ID:      uint64(row.BuildingID),
			Name:    row.BuildingName,
			Address: row.BuildingAddress,
		},
	}
}

func mapRoomFromListByBuilding(row db.ListRoomsByBuildingRow) *entity.Room {
	return &entity.Room{
		ID:         uint64(row.ID),
		Name:       row.Name,
		BuildingID: uint64(row.BuildingID),
		Building: &entity.Building{
			ID:      uint64(row.BuildingID),
			Name:    row.BuildingName,
			Address: row.BuildingAddress,
		},
	}
}

func mapRoomFromGetByNameAndBuilding(row db.GetRoomByNameAndBuildingRow) *entity.Room {
	return &entity.Room{
		ID:         uint64(row.ID),
		Name:       row.Name,
		BuildingID: uint64(row.BuildingID),
		Building: &entity.Building{
			ID:      uint64(row.BuildingID),
			Name:    row.BuildingName,
			Address: row.BuildingAddress,
		},
	}
}
