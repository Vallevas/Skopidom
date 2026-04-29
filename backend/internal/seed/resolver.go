// Package seed provides functionality for importing inventory items from CSV files.
package seed

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/pkg/logger"
)

// Resolver handles finding or creating reference data (categories, buildings, rooms).
type Resolver struct {
	categories repository.CategoryRepository
	buildings  repository.BuildingRepository
	rooms      repository.RoomRepository
}

// NewResolver creates a new Resolver.
func NewResolver(
	categories repository.CategoryRepository,
	buildings repository.BuildingRepository,
	rooms repository.RoomRepository,
) *Resolver {
	return &Resolver{
		categories: categories,
		buildings:  buildings,
		rooms:      rooms,
	}
}

// ResolveCategory finds a category by name or creates it if it doesn't exist.
func (r *Resolver) ResolveCategory(ctx context.Context, name string) (uint64, error) {
	// Try to find existing category
	cat, err := r.categories.GetByName(ctx, name)
	if err == nil {
		return cat.ID, nil
	}

	// If not found, create new category
	if errors.Is(err, logger.ErrNotFound) {
		newCat := &entity.Category{Name: name}
		if err := r.categories.Create(ctx, newCat); err != nil {
			return 0, fmt.Errorf("create category %q: %w", name, err)
		}
		return newCat.ID, nil
	}

	return 0, fmt.Errorf("get category by name: %w", err)
}

// ResolveBuilding finds a building by name or creates it if it doesn't exist.
func (r *Resolver) ResolveBuilding(ctx context.Context, name string) (uint64, error) {
	// Try to find existing building
	building, err := r.buildings.GetByName(ctx, name)
	if err == nil {
		return building.ID, nil
	}

	// If not found, create new building with empty address
	if errors.Is(err, logger.ErrNotFound) {
		newBuilding := &entity.Building{
			Name:    name,
			Address: "", // Empty address by default
		}
		if err := r.buildings.Create(ctx, newBuilding); err != nil {
			return 0, fmt.Errorf("create building %q: %w", name, err)
		}
		return newBuilding.ID, nil
	}

	return 0, fmt.Errorf("get building by name: %w", err)
}

// ResolveRoom finds a room by name and building ID or creates it if it doesn't exist.
func (r *Resolver) ResolveRoom(ctx context.Context, name string, buildingID uint64) (uint64, error) {
	// Try to find existing room
	room, err := r.rooms.GetByNameAndBuilding(ctx, name, buildingID)
	if err == nil {
		return room.ID, nil
	}

	// If not found, create new room
	if errors.Is(err, logger.ErrNotFound) {
		newRoom := &entity.Room{
			Name:       name,
			BuildingID: buildingID,
		}
		if err := r.rooms.Create(ctx, newRoom); err != nil {
			return 0, fmt.Errorf("create room %q in building %d: %w", name, buildingID, err)
		}
		return newRoom.ID, nil
	}

	return 0, fmt.Errorf("get room by name and building: %w", err)
}
