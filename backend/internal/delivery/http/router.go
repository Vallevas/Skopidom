// Package http provides the HTTP server setup for the inventory API.
package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vallevas/Skopidom/internal/delivery/http/handler"
	"github.com/Vallevas/Skopidom/internal/delivery/http/middleware"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/internal/infrastructure/storage"
	itemUC "github.com/Vallevas/Skopidom/internal/usecase/item"
	userUC "github.com/Vallevas/Skopidom/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// RouterConfig holds all dependencies required to build the HTTP router.
type RouterConfig struct {
	ItemUC  itemUC.UseCase
	UserUC  userUC.UseCase
	Repos   RepoSet
	Storage storage.FileStorage

	JWTSecret string
	JWTTTL    time.Duration

	// StaticDir is the filesystem directory served at /static/.
	StaticDir string

	// DevMode enables verbose error responses (detail field in JSON).
	DevMode bool
	// AllowedOrigins is passed directly to the CORS middleware.
	// Use []string{"*"} for development and specific domains for production.
	AllowedOrigins []string
}

// RepoSet groups lookup repositories needed to construct simple use cases.
type RepoSet struct {
	Categories repository.CategoryRepository
	Buildings  repository.BuildingRepository
	Rooms      repository.RoomRepository
}

// NewRouter constructs and returns a fully-wired chi.Router.
func NewRouter(cfg RouterConfig) http.Handler {
	// Initialise error verbosity for all handlers once at startup.
	handler.InitErrorMode(cfg.DevMode)

	r := chi.NewRouter()

	// Global middleware stack.
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// Serve uploaded photos as static files.
	if cfg.StaticDir != "" {
		r.Handle("/static/*",
			http.StripPrefix("/static/",
				http.FileServer(http.Dir(cfg.StaticDir))))
	}

	// Instantiate simple use cases for lookup entities.
	catUC := handler.NewSimpleCategoryUC(cfg.Repos.Categories)
	buildingUC := handler.NewSimpleBuildingUC(cfg.Repos.Buildings)
	roomUC := handler.NewSimpleRoomUC(cfg.Repos.Rooms, cfg.Repos.Buildings)

	// Instantiate handlers.
	itemH := handler.NewItemHandler(cfg.ItemUC)
	userH := handler.NewUserHandler(cfg.UserUC)
	catH := handler.NewCategoryHandler(catUC)
	buildingH := handler.NewBuildingHandler(buildingUC)
	roomH := handler.NewRoomHandler(roomUC)
	uploadH := handler.NewUploadHandler(cfg.Storage, cfg.ItemUC)
	loginH := newLoginHandler(cfg.UserUC, cfg.JWTSecret, cfg.JWTTTL)

	authMW := middleware.Auth(cfg.JWTSecret)

	r.Route("/api/v1", func(r chi.Router) {

		// Public: login only.
		r.Post("/auth/login", loginH)

		// All routes below require a valid JWT.
		r.Group(func(r chi.Router) {
			r.Use(authMW)

			// ── Items ──────────────────────────────────────────────────────
			r.Route("/items", func(r chi.Router) {
				r.Get("/", itemH.List)
				r.Post("/", itemH.Create)
				r.Get("/barcode/{barcode}", itemH.GetByBarcode)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", itemH.GetByID)
					r.Get("/audit", itemH.GetAuditLog)
					r.Get("/photos", uploadH.ListItemPhotos)
					r.Patch("/", itemH.Update)
					r.Patch("/room", itemH.MoveToRoom)
					r.Post("/photos", uploadH.UploadItemPhoto)
					r.Delete("/photos/{photo_id}", uploadH.DeleteItemPhoto)
					r.With(middleware.RequireAdmin).Delete("/", itemH.Dispose)
				})
			})

			// ── Lookup / reference data (read: any auth; write: admin) ─────
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", catH.List)
				r.With(middleware.RequireAdmin).Post("/", catH.Create)
				r.With(middleware.RequireAdmin).Patch("/{id}", catH.Update)
				r.With(middleware.RequireAdmin).Delete("/{id}", catH.Delete)
			})

			r.Route("/buildings", func(r chi.Router) {
				r.Get("/", buildingH.List)
				r.With(middleware.RequireAdmin).Post("/", buildingH.Create)
				r.With(middleware.RequireAdmin).Patch("/{id}", buildingH.Update)
				r.With(middleware.RequireAdmin).Delete("/{id}", buildingH.Delete)
			})

			r.Route("/rooms", func(r chi.Router) {
				r.Get("/", roomH.List)
				r.With(middleware.RequireAdmin).Post("/", roomH.Create)
				r.With(middleware.RequireAdmin).Patch("/{id}", roomH.Update)
				r.With(middleware.RequireAdmin).Delete("/{id}", roomH.Delete)
			})

			// ── User management (admin only) ───────────────────────────────
			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RequireAdmin)
				r.Get("/", userH.List)
				r.Post("/", userH.Register)
				r.Get("/{id}", userH.GetByID)
				r.Patch("/{id}", userH.Update)
				r.Delete("/{id}", userH.Delete)
			})
		})
	})

	return r
}

// newLoginHandler returns an http.HandlerFunc for POST /api/v1/auth/login.
func newLoginHandler(
	uc userUC.UseCase,
	secret string,
	ttl time.Duration,
) http.HandlerFunc {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			routerJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := uc.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			routerJSONError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := middleware.GenerateToken(user.ID, user.Role, secret, ttl)
		if err != nil {
			routerJSONError(w, http.StatusInternalServerError, "could not generate token")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"token": token,
			"user":  user,
		})
	}
}

// routerJSONError writes a simple JSON error response from within this package.
func routerJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + msg + `"}`))
}
