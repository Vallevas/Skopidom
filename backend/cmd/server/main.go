// Package main is the entry point for the inventory HTTP server.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryHTTP "github.com/Vallevas/Skopidom/internal/delivery/http"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres"
	"github.com/Vallevas/Skopidom/internal/infrastructure/storage"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	userUC "github.com/Vallevas/Skopidom/internal/usecase/user"
	"github.com/Vallevas/Skopidom/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if cfg.Debug {
		log.Println("running in DEBUG mode — verbose errors enabled")
	}

	// ── Database ───────────────────────────────────────────────────────────
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	// ── Migrations ─────────────────────────────────────────────────────────
	if err := postgres.RunMigrations(cfg.Postgres.DSN, cfg.Postgres.MigrationsPath); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations applied")

	// ── Repositories ───────────────────────────────────────────────────────────
	itemRepo := postgres.NewItemRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	categoryRepo := postgres.NewCategoryRepo(pool)
	buildingRepo := postgres.NewBuildingRepo(pool)
	roomRepo := postgres.NewRoomRepo(pool)

	// ── Audit Logger ───────────────────────────────────────────────────────────
	//auditLogger :=  blockchain.NewBlockchainAuditLogger(...)
	auditLogger := postgres.NewPostgresAuditLogger(pool)

	// ── Use Cases ──────────────────────────────────────────────────────────────
	itemUseCase := itemUC.New(itemRepo, categoryRepo, roomRepo, auditLogger)
	userUseCase := userUC.New(userRepo)

	// ── File Storage ───────────────────────────────────────────────────────
	fileStorage, err := storage.NewLocalStorage(
		cfg.Storage.Dir,
		cfg.Storage.BaseURL,
	)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	// ── HTTP Router ────────────────────────────────────────────────────────
	router := deliveryHTTP.NewRouter(deliveryHTTP.RouterConfig{
		ItemUC:  itemUseCase,
		UserUC:  userUseCase,
		Storage: fileStorage,
		Repos: deliveryHTTP.RepoSet{
			Categories: categoryRepo,
			Buildings:  buildingRepo,
			Rooms:      roomRepo,
		},
		JWTSecret:      cfg.JWT.Secret,
		JWTTTL:         cfg.JWT.TTL,
		StaticDir:      cfg.Storage.Dir,
		DevMode:        cfg.IsDevelopment(),
		AllowedOrigins: cfg.Server.AllowedOrigins,
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// ── Graceful shutdown ──────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on :%s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
