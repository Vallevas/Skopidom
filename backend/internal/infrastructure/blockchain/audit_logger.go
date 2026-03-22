// Package blockchain will contain the Ethereum smart contract integration.
// This file defines BlockchainAuditLogger — the future implementation of
// repository.AuditLogger that anchors audit events on-chain.
//
// To activate: replace postgres.NewPostgresAuditLogger in main.go with
// blockchain.NewBlockchainAuditLogger(pool, contractAddress, ethClient).
// No other code changes are required.
package blockchain

import (
	"context"
	"log/slog"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"
	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// BlockchainAuditLogger implements repository.AuditLogger.
// Every event is first persisted to PostgreSQL (via the embedded postgres
// logger), then anchored to the Ethereum chain. The tx_hash is written back
// once the transaction is confirmed.
//
// If the chain is unreachable the event remains in PostgreSQL without a
// tx_hash — no data is lost and the primary operation is not affected.
// A background reconciliation job (not yet implemented) can retry
// unanchored events identified by WHERE tx_hash = ”.
type BlockchainAuditLogger struct {
	db      repository.AuditLogger // PostgreSQL fallback — always written first.
	queries *db.Queries            // needed to update tx_hash after confirmation.

	// contractAddress is the deployed InventoryAudit contract address.
	// contractAddress string  (uncomment when go-ethereum bindings are added)

	// ethClient is the go-ethereum RPC client.
	// ethClient *ethclient.Client  (uncomment when go-ethereum is added)
}

// NewBlockchainAuditLogger constructs a BlockchainAuditLogger.
// contractAddress and ethClient parameters will be added once go-ethereum
// bindings are generated from the Solidity contract.
func NewBlockchainAuditLogger(
	pool *pgxpool.Pool,
	postgresLogger repository.AuditLogger,
) *BlockchainAuditLogger {
	return &BlockchainAuditLogger{
		db:      postgresLogger,
		queries: db.New(stdlib.OpenDBFromPool(pool)),
	}
}

// Log persists the event to PostgreSQL then submits an Ethereum transaction.
func (l *BlockchainAuditLogger) Log(ctx context.Context, event *entity.AuditEvent) error {
	// Step 1: Always write to PostgreSQL first — data is safe even if chain fails.
	if err := l.db.Log(ctx, event); err != nil {
		return err
	}

	// Step 2: Submit transaction to Ethereum.
	// TODO: uncomment once go-ethereum bindings are generated.
	//
	// txHash, err := l.contract.LogEvent(&bind.TransactOpts{...}, event.ItemID, event.Action, event.Payload)
	// if err != nil {
	//     slog.Warn("audit: blockchain submission failed — event stored in postgres only",
	//         "audit_event_id", event.ID, "err", err)
	//     return nil  // fallback: postgres record exists without tx_hash
	// }
	//
	// Step 3: Write tx_hash back to the postgres record.
	// if err := l.queries.UpdateAuditEventTxHash(ctx, db.UpdateAuditEventTxHashParams{
	//     TxHash: txHash.Hex(),
	//     ID:     int64(event.ID),
	// }); err != nil {
	//     slog.Error("audit: failed to persist tx_hash", "err", err)
	// }
	// event.TxHash = txHash.Hex()

	slog.Warn("audit: blockchain not yet connected — event stored in postgres only",
		"audit_event_id", event.ID)
	return nil
}

// ListByItem delegates to the PostgreSQL implementation.
func (l *BlockchainAuditLogger) ListByItem(
	ctx context.Context,
	itemID uint64,
) ([]*entity.AuditEvent, error) {
	return l.db.ListByItem(ctx, itemID)
}
