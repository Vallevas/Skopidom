// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// PostgresAuditLogger implements repository.AuditLogger using PostgreSQL.
// It is the default (non-blockchain) implementation — all events are stored
// in the audit_events table with an empty tx_hash.
//
// To enable blockchain anchoring, wrap this struct inside a
// BlockchainAuditLogger that calls Log here first, then submits the
// transaction and calls UpdateAuditEventTxHash.
type PostgresAuditLogger struct {
	queries *db.Queries
}

// NewPostgresAuditLogger constructs a PostgresAuditLogger backed by the pool.
func NewPostgresAuditLogger(pool *pgxpool.Pool) *PostgresAuditLogger {
	return &PostgresAuditLogger{queries: db.New(stdlib.OpenDBFromPool(pool))}
}

// Log persists an audit event to PostgreSQL.
// Failures are logged via slog but never propagate to the caller — a logging
// failure must not roll back the primary business operation.
func (l *PostgresAuditLogger) Log(ctx context.Context, event *entity.AuditEvent) error {
	row, err := l.queries.CreateAuditEvent(ctx, db.CreateAuditEventParams{
		ItemID:  int64(event.ItemID),
		ActorID: int64(event.ActorID),
		Action:  string(event.Action),
		// string → json.RawMessage: JSONB column expects raw JSON bytes.
		Payload: json.RawMessage(event.Payload),
	})
	if err != nil {
		// Log internally — do not propagate.
		slog.Error("audit: failed to persist event",
			"item_id", event.ItemID,
			"action", event.Action,
			"err", err,
		)
		return nil
	}

	event.ID = uint64(row.ID)
	event.CreatedAt = row.CreatedAt
	return nil
}

// ListByItem returns all audit events for the given item ordered by time.
func (l *PostgresAuditLogger) ListByItem(
	ctx context.Context,
	itemID uint64,
) ([]*entity.AuditEvent, error) {
	rows, err := l.queries.ListAuditEventsByItem(ctx, int64(itemID))
	if err != nil {
		return nil, fmt.Errorf("AuditLogger.ListByItem: %w", err)
	}

	events := make([]*entity.AuditEvent, len(rows))
	for i, row := range rows {
		events[i] = mapAuditEvent(row)
	}
	return events, nil
}

// ── mapping ───────────────────────────────────────────────────────────────────

func mapAuditEvent(row db.ListAuditEventsByItemRow) *entity.AuditEvent {
	return &entity.AuditEvent{
		ID:      uint64(row.ID),
		ItemID:  uint64(row.ItemID),
		ActorID: uint64(row.ActorID),
		Actor: &entity.User{
			ID:       uint64(row.ActorID),
			FullName: row.ActorFullName,
			Email:    row.ActorEmail,
			Role:     entity.UserRole(row.ActorRole),
		},
		Action: entity.AuditAction(row.Action),
		// json.RawMessage → string: domain entity uses plain string for portability.
		Payload:   string(row.Payload),
		TxHash:    row.TxHash,
		CreatedAt: row.CreatedAt,
	}
}

