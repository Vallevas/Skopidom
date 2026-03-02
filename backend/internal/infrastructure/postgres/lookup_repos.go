// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── CategoryRepo ─────────────────────────────────────────────────────────────

// CategoryRepo implements repository.CategoryRepository using PostgreSQL.
type CategoryRepo struct {
	pool *pgxpool.Pool
}

// NewCategoryRepo constructs a CategoryRepo backed by the given pool.
func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{pool: pool}
}

func (r *CategoryRepo) Create(ctx context.Context, cat *entity.Category) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO categories (name) VALUES ($1) RETURNING id`,
		cat.Name,
	).Scan(&cat.ID)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Create: %w", err)
	}
	return nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uint64) (*entity.Category, error) {
	cat := &entity.Category{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name FROM categories WHERE id = $1`, id,
	).Scan(&cat.ID, &cat.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.GetByID: %w", err)
	}
	return cat, nil
}

func (r *CategoryRepo) List(ctx context.Context) ([]*entity.Category, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name FROM categories ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("CategoryRepo.List: %w", err)
	}
	defer rows.Close()

	cats := make([]*entity.Category, 0)
	for rows.Next() {
		cat := &entity.Category{}
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("CategoryRepo.List scan: %w", err)
		}
		cats = append(cats, cat)
	}
	return cats, rows.Err()
}

func (r *CategoryRepo) Update(ctx context.Context, cat *entity.Category) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE categories SET name = $1 WHERE id = $2`,
		cat.Name, cat.ID,
	)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id uint64) error {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM categories WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("CategoryRepo.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ── BuildingRepo ──────────────────────────────────────────────────────────────

// BuildingRepo implements repository.BuildingRepository using PostgreSQL.
type BuildingRepo struct {
	pool *pgxpool.Pool
}

// NewBuildingRepo constructs a BuildingRepo backed by the given pool.
func NewBuildingRepo(pool *pgxpool.Pool) *BuildingRepo {
	return &BuildingRepo{pool: pool}
}

func (r *BuildingRepo) Create(ctx context.Context, b *entity.Building) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO buildings (name, address) VALUES ($1, $2) RETURNING id`,
		b.Name, b.Address,
	).Scan(&b.ID)
	if err != nil {
		return fmt.Errorf("BuildingRepo.Create: %w", err)
	}
	return nil
}

func (r *BuildingRepo) GetByID(ctx context.Context, id uint64) (*entity.Building, error) {
	b := &entity.Building{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, address FROM buildings WHERE id = $1`, id,
	).Scan(&b.ID, &b.Name, &b.Address)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("BuildingRepo.GetByID: %w", err)
	}
	return b, nil
}

func (r *BuildingRepo) List(ctx context.Context) ([]*entity.Building, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, address FROM buildings ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("BuildingRepo.List: %w", err)
	}
	defer rows.Close()

	buildings := make([]*entity.Building, 0)
	for rows.Next() {
		b := &entity.Building{}
		if err := rows.Scan(&b.ID, &b.Name, &b.Address); err != nil {
			return nil, fmt.Errorf("BuildingRepo.List scan: %w", err)
		}
		buildings = append(buildings, b)
	}
	return buildings, rows.Err()
}

func (r *BuildingRepo) Update(ctx context.Context, b *entity.Building) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE buildings SET name = $1, address = $2 WHERE id = $3`,
		b.Name, b.Address, b.ID,
	)
	if err != nil {
		return fmt.Errorf("BuildingRepo.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *BuildingRepo) Delete(ctx context.Context, id uint64) error {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM buildings WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("BuildingRepo.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// ── RoomRepo ─────────────────────────────────────────────────────────────────

// RoomRepo implements repository.RoomRepository using PostgreSQL.
type RoomRepo struct {
	pool *pgxpool.Pool
}

// NewRoomRepo constructs a RoomRepo backed by the given pool.
func NewRoomRepo(pool *pgxpool.Pool) *RoomRepo {
	return &RoomRepo{pool: pool}
}

func (r *RoomRepo) Create(ctx context.Context, room *entity.Room) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO rooms (name, building_id) VALUES ($1, $2) RETURNING id`,
		room.Name, room.BuildingID,
	).Scan(&room.ID)
	if err != nil {
		return fmt.Errorf("RoomRepo.Create: %w", err)
	}
	return nil
}

func (r *RoomRepo) GetByID(ctx context.Context, id uint64) (*entity.Room, error) {
	room := &entity.Room{Building: &entity.Building{}}
	err := r.pool.QueryRow(ctx,
		`SELECT r.id, r.name, r.building_id, b.name, b.address
		 FROM rooms r
		 JOIN buildings b ON b.id = r.building_id
		 WHERE r.id = $1`, id,
	).Scan(&room.ID, &room.Name, &room.BuildingID,
		&room.Building.Name, &room.Building.Address)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.GetByID: %w", err)
	}
	room.Building.ID = room.BuildingID
	return room, nil
}

func (r *RoomRepo) List(ctx context.Context) ([]*entity.Room, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.name, r.building_id, b.name, b.address
		 FROM rooms r
		 JOIN buildings b ON b.id = r.building_id
		 ORDER BY b.name, r.name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.List: %w", err)
	}
	defer rows.Close()

	rooms := make([]*entity.Room, 0)
	for rows.Next() {
		room := &entity.Room{Building: &entity.Building{}}
		if err := rows.Scan(
			&room.ID, &room.Name, &room.BuildingID,
			&room.Building.Name, &room.Building.Address,
		); err != nil {
			return nil, fmt.Errorf("RoomRepo.List scan: %w", err)
		}
		room.Building.ID = room.BuildingID
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *RoomRepo) ListByBuilding(
	ctx context.Context,
	buildingID uint64,
) ([]*entity.Room, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.name, r.building_id, b.name, b.address
		 FROM rooms r
		 JOIN buildings b ON b.id = r.building_id
		 WHERE r.building_id = $1
		 ORDER BY r.name ASC`,
		buildingID,
	)
	if err != nil {
		return nil, fmt.Errorf("RoomRepo.ListByBuilding: %w", err)
	}
	defer rows.Close()

	rooms := make([]*entity.Room, 0)
	for rows.Next() {
		room := &entity.Room{Building: &entity.Building{}}
		if err := rows.Scan(
			&room.ID, &room.Name, &room.BuildingID,
			&room.Building.Name, &room.Building.Address,
		); err != nil {
			return nil, fmt.Errorf("RoomRepo.ListByBuilding scan: %w", err)
		}
		room.Building.ID = room.BuildingID
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (r *RoomRepo) Update(ctx context.Context, room *entity.Room) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE rooms SET name = $1, building_id = $2 WHERE id = $3`,
		room.Name, room.BuildingID, room.ID,
	)
	if err != nil {
		return fmt.Errorf("RoomRepo.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *RoomRepo) Delete(ctx context.Context, id uint64) error {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM rooms WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("RoomRepo.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
