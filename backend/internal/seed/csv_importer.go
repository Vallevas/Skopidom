// Package seed provides functionality for importing inventory items from CSV files.
package seed

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Vallevas/Skopidom/internal/domain/repository"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
)

// CSVRow represents a single row from the CSV file.
type CSVRow struct {
	Barcode         string
	InventoryNumber string
	Name            string
	Category        string
	Building        string
	Room            string
	Description     string
}

// ImportResult represents the result of importing a single row.
type ImportResult struct {
	Row     CSVRow
	Success bool
	Skipped bool
	Error   error
}

// Importer handles CSV import operations.
type Importer struct {
	itemUC   itemUC.UseCase
	resolver *Resolver
	itemRepo repository.ItemRepository
	userID   uint64
}

// NewImporter creates a new CSV importer.
func NewImporter(
	itemUC itemUC.UseCase,
	categories repository.CategoryRepository,
	buildings repository.BuildingRepository,
	rooms repository.RoomRepository,
	itemRepo repository.ItemRepository,
	userID uint64,
) *Importer {
	return &Importer{
		itemUC:   itemUC,
		resolver: NewResolver(categories, buildings, rooms),
		itemRepo: itemRepo,
		userID:   userID,
	}
}

// ParseCSV reads and parses a CSV file into CSVRow structs.
func ParseCSV(filePath string) ([]CSVRow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	// Validate header
	expectedHeader := []string{"barcode", "inventory_number", "name", "category", "building", "room", "description"}
	if len(header) != len(expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header: expected %d columns, got %d", len(expectedHeader), len(header))
	}

	var rows []CSVRow
	lineNum := 1 // header is line 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV line %d: %w", lineNum+1, err)
		}
		lineNum++

		if len(record) != len(expectedHeader) {
			return nil, fmt.Errorf("line %d: expected %d columns, got %d", lineNum, len(expectedHeader), len(record))
		}

		rows = append(rows, CSVRow{
			Barcode:         strings.TrimSpace(record[0]),
			InventoryNumber: strings.TrimSpace(record[1]),
			Name:            strings.TrimSpace(record[2]),
			Category:        strings.TrimSpace(record[3]),
			Building:        strings.TrimSpace(record[4]),
			Room:            strings.TrimSpace(record[5]),
			Description:     strings.TrimSpace(record[6]),
		})
	}

	return rows, nil
}

// ValidateRow checks if a CSV row has all required fields.
func ValidateRow(row CSVRow) error {
	if row.Barcode == "" {
		return fmt.Errorf("barcode is required")
	}
	if row.InventoryNumber == "" {
		return fmt.Errorf("inventory_number is required")
	}
	if row.Name == "" {
		return fmt.Errorf("name is required")
	}
	if row.Category == "" {
		return fmt.Errorf("category is required")
	}
	if row.Building == "" {
		return fmt.Errorf("building is required")
	}
	if row.Room == "" {
		return fmt.Errorf("room is required")
	}
	return nil
}

// ImportRow imports a single CSV row into the database.
func (imp *Importer) ImportRow(ctx context.Context, row CSVRow) ImportResult {
	result := ImportResult{Row: row}

	// Validate row
	if err := ValidateRow(row); err != nil {
		result.Error = fmt.Errorf("validation failed: %w", err)
		return result
	}

	// Check if barcode already exists
	barcodeExists, err := imp.itemRepo.BarcodeExists(ctx, row.Barcode)
	if err != nil {
		result.Error = fmt.Errorf("check barcode exists: %w", err)
		return result
	}
	if barcodeExists {
		result.Skipped = true
		result.Error = fmt.Errorf("barcode %q already exists", row.Barcode)
		return result
	}

	// Check if inventory number already exists
	invExists, err := imp.itemRepo.InventoryNumberExists(ctx, row.InventoryNumber)
	if err != nil {
		result.Error = fmt.Errorf("check inventory_number exists: %w", err)
		return result
	}
	if invExists {
		result.Skipped = true
		result.Error = fmt.Errorf("inventory_number %q already exists", row.InventoryNumber)
		return result
	}

	// Resolve category (find or create)
	categoryID, err := imp.resolver.ResolveCategory(ctx, row.Category)
	if err != nil {
		result.Error = fmt.Errorf("resolve category: %w", err)
		return result
	}

	// Resolve building (find or create)
	buildingID, err := imp.resolver.ResolveBuilding(ctx, row.Building)
	if err != nil {
		result.Error = fmt.Errorf("resolve building: %w", err)
		return result
	}

	// Resolve room (find or create)
	roomID, err := imp.resolver.ResolveRoom(ctx, row.Room, buildingID)
	if err != nil {
		result.Error = fmt.Errorf("resolve room: %w", err)
		return result
	}

	// Create item using the use case (this will also create audit log)
	input := itemUC.CreateInput{
		Barcode:         row.Barcode,
		InventoryNumber: row.InventoryNumber,
		Name:            row.Name,
		CategoryID:      categoryID,
		RoomID:          roomID,
		Description:     row.Description,
		ActorID:         imp.userID,
	}

	_, err = imp.itemUC.Create(ctx, input)
	if err != nil {
		result.Error = fmt.Errorf("create item: %w", err)
		return result
	}

	result.Success = true
	return result
}
