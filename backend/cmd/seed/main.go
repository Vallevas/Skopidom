// Package main provides a CLI tool for importing inventory items from CSV files.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres"
	"github.com/Vallevas/Skopidom/internal/seed"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	"github.com/Vallevas/Skopidom/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command-line flags
	var (
		filePath  = flag.String("file", "", "Path to CSV file (required)")
		userEmail = flag.String("user", "", "Email of the user performing the import (required)")
		dryRun    = flag.Bool("dry-run", false, "Parse and validate CSV without importing")
	)
	flag.Parse()

	// Validate required flags
	if *filePath == "" {
		log.Fatal("Error: --file flag is required\nUsage: go run ./cmd/seed --file=items.csv --user=admin@university.ru")
	}
	if *userEmail == "" {
		log.Fatal("Error: --user flag is required\nUsage: go run ./cmd/seed --file=items.csv --user=admin@university.ru")
	}

	// Load environment variables
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse CSV file
	fmt.Printf("Parsing CSV file: %s\n", *filePath)
	rows, err := seed.ParseCSV(*filePath)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}
	fmt.Printf("Parsed %d rows from CSV\n\n", len(rows))

	// Validate all rows
	var validationErrors int
	for i, row := range rows {
		if err := seed.ValidateRow(row); err != nil {
			fmt.Printf("[%d/%d] ✗ Validation error: %v (barcode: %s)\n", i+1, len(rows), err, row.Barcode)
			validationErrors++
		}
	}

	if validationErrors > 0 {
		log.Fatalf("\n%d validation errors found. Please fix the CSV file and try again.", validationErrors)
	}

	fmt.Println("✓ All rows passed validation")

	// If dry-run, exit here
	if *dryRun {
		fmt.Println("\n--dry-run mode: No data was imported")
		return
	}

	fmt.Println()

	// Connect to database
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	itemRepo := postgres.NewItemRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	categoryRepo := postgres.NewCategoryRepo(pool)
	buildingRepo := postgres.NewBuildingRepo(pool)
	roomRepo := postgres.NewRoomRepo(pool)
	photoRepo := postgres.NewPhotoRepo(pool)
	disposalDocRepo := postgres.NewDisposalDocumentRepo(pool)
	auditLogger := postgres.NewPostgresAuditLogger(pool)

	// Find user by email
	user, err := userRepo.GetByEmail(ctx, *userEmail)
	if err != nil {
		log.Fatalf("Failed to find user with email %q: %v", *userEmail, err)
	}
	fmt.Printf("Importing as: %s (%s) [ID: %d]\n\n", user.FullName, user.Email, user.ID)

	// Initialize use case
	itemUseCase := itemUC.New(itemRepo, categoryRepo, roomRepo, photoRepo, disposalDocRepo, auditLogger)

	// Initialize importer
	importer := seed.NewImporter(itemUseCase, categoryRepo, buildingRepo, roomRepo, itemRepo, user.ID)

	// Import rows
	fmt.Printf("Importing %d items...\n\n", len(rows))
	startTime := time.Now()

	var (
		successCount int
		skipCount    int
		errorCount   int
	)

	for i, row := range rows {
		result := importer.ImportRow(ctx, row)

		if result.Success {
			fmt.Printf("[%d/%d] ✓ Created: %s (barcode: %s)\n", i+1, len(rows), row.Name, row.Barcode)
			successCount++
		} else if result.Skipped {
			fmt.Printf("[%d/%d] ⊘ Skipped: %v\n", i+1, len(rows), result.Error)
			skipCount++
		} else {
			fmt.Printf("[%d/%d] ✗ Error: %v (barcode: %s)\n", i+1, len(rows), result.Error, row.Barcode)
			errorCount++
		}
	}

	duration := time.Since(startTime)

	// Print summary
	fmt.Println("\n" + strings.Repeat("─", 50))
	fmt.Println("Import Summary:")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("  ✓ Created:  %d\n", successCount)
	fmt.Printf("  ⊘ Skipped:  %d\n", skipCount)
	fmt.Printf("  ✗ Errors:   %d\n", errorCount)
	fmt.Printf("  Total:      %d\n", len(rows))
	fmt.Printf("  Duration:   %s\n", duration.Round(time.Millisecond))
	fmt.Println(strings.Repeat("─", 50))

	if errorCount > 0 {
		os.Exit(1)
	}
}
