// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Vallevas/Skopidom/pkg/logger"
)

// respond serialises v as JSON and writes it with the given HTTP status code.
func respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// decodeJSON reads and validates a JSON request body into dst.
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// errorResponse is the standard JSON error envelope.
type errorResponse struct {
	Error string `json:"error"`
}

// handleError maps domain/app errors to appropriate HTTP status codes.
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		respond(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, apperrors.ErrAlreadyExists):
		respond(w, http.StatusConflict, errorResponse{Error: err.Error()})
	case errors.Is(err, apperrors.ErrDisposed):
		respond(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
	case errors.Is(err, apperrors.ErrForbidden):
		respond(w, http.StatusForbidden, errorResponse{Error: err.Error()})
	case errors.Is(err, apperrors.ErrUnauthorized):
		respond(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
	case errors.Is(err, apperrors.ErrInvalidInput):
		respond(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		respond(w, http.StatusInternalServerError,
			errorResponse{Error: "internal server error"})
	}
}
