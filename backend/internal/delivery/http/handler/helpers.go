// Package handler provides HTTP handler implementations for the inventory API.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
// Detail is only populated in debug mode.
type errorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

// isDev controls whether full error details are included in responses.
// Set once at startup via InitErrorMode — avoids threading Config through
// every handler constructor.
var isDev bool

// InitErrorMode sets the error verbosity for all handlers in this package.
// Must be called once during application startup before serving requests.
func InitErrorMode(debug bool) {
	isDev = debug
}

// handleError maps domain/app errors to appropriate HTTP status codes.
//
// Safe messages (ErrNotFound, ErrForbidden, etc.) are always shown to the
// client verbatim. Internal errors are hidden in production — only the
// generic "internal server error" string is returned — but are logged via
// slog regardless of mode so that operators can investigate.
//
// In debug mode the full error chain is additionally included in "detail".
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, logger.ErrNotFound):
		respond(w, http.StatusNotFound,
			userFacingError("resource not found", err))

	case errors.Is(err, logger.ErrAlreadyExists):
		respond(w, http.StatusConflict,
			userFacingError("resource already exists", err))

	case errors.Is(err, logger.ErrDisposed):
		respond(w, http.StatusUnprocessableEntity,
			userFacingError("item is disposed and cannot be modified", err))

	case errors.Is(err, logger.ErrForbidden):
		respond(w, http.StatusForbidden,
			userFacingError("insufficient permissions", err))

	case errors.Is(err, logger.ErrUnauthorized):
		respond(w, http.StatusUnauthorized,
			userFacingError("unauthorized", err))

	case errors.Is(err, logger.ErrInvalidInput):
		// Validation messages describe the request — always safe to expose.
		respond(w, http.StatusBadRequest, errorResponse{Error: err.Error()})

	case errors.Is(err, logger.ErrConflict):
		respond(w, http.StatusConflict,
			userFacingError(err.Error(), err))

	default:
		// Always log the full error for operators.
		slog.Error("unhandled internal error", "err", err)

		resp := errorResponse{Error: "internal server error"}
		if isDev {
			resp.Detail = err.Error()
		}
		respond(w, http.StatusInternalServerError, resp)
	}
}

// userFacingError builds an errorResponse with a safe public message.
// In debug mode the full error chain is appended to Detail for easier
// debugging without exposing internals to production users.
func userFacingError(publicMsg string, err error) errorResponse {
	resp := errorResponse{Error: publicMsg}
	if isDev {
		resp.Detail = err.Error()
	}
	return resp
}

// wrapInvalidInput wraps a generic decode error as ErrInvalidInput.
func wrapInvalidInput(err error) error {
	return fmt.Errorf("%s: %w", err.Error(), logger.ErrInvalidInput)
}
