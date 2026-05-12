// Package blockchain contains the Ethereum smart contract integration.
package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/Vallevas/Skopidom/internal/domain/entity"
	"github.com/Vallevas/Skopidom/internal/domain/repository"

	"github.com/Vallevas/Skopidom/internal/infrastructure/postgres"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate abigen --abi ../../blockchain/build/InventoryAudit.abi --pkg blockchain --type InventoryAudit --out ./inventory_audit.go

// BlockchainAuditLogger implements repository.AuditLogger.
type BlockchainAuditLogger struct {
	db              repository.AuditLogger
	contract        *Blockchain // Changed from *InventoryAudit to *Blockchain
	contractAddress common.Address
	ethClient       *ethclient.Client
	transactOpts    *bind.TransactOpts
}

// NewBlockchainAuditLogger constructs a BlockchainAuditLogger.
func NewBlockchainAuditLogger(
	pool *pgxpool.Pool,
	contractAddress string,
	ethClient *ethclient.Client,
	privateKey string,
) (*BlockchainAuditLogger, error) {
	// Parse contract address
	addr := common.HexToAddress(contractAddress)
	if addr == (common.Address{}) {
		return nil, fmt.Errorf("invalid contract address: %s", contractAddress)
	}

	// Create transact opts from private key
	transactOpts, err := createTransactOpts(privateKey, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create transact opts: %w", err)
	}

	// Bind to deployed contract - Use NewBlockchain instead of NewInventoryAudit
	contract, err := NewBlockchain(addr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract: %w", err)
	}

	return &BlockchainAuditLogger{
		db:              postgres.NewPostgresAuditLogger(pool),
		contract:        contract,
		contractAddress: addr,
		ethClient:       ethClient,
		transactOpts:    transactOpts,
	}, nil
}

// createTransactOpts creates transaction options from a private key
func createTransactOpts(privateKeyHex string, client *ethclient.Client) (*bind.TransactOpts, error) {
	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	// Parse private key
	key, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transact opts
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	// Set gas limit (adjust as needed)
	auth.GasLimit = uint64(300000)

	return auth, nil
}

// Log persists the event to PostgreSQL then submits an Ethereum transaction.
func (l *BlockchainAuditLogger) Log(ctx context.Context, event *entity.AuditEvent) error {
	// Step 1: Always write to PostgreSQL first
	if err := l.db.Log(ctx, event); err != nil {
		return err
	}

	// Step 2: Submit transaction to Ethereum
	txHash, err := l.logToBlockchain(ctx, event)
	if err != nil {
		slog.Warn("audit: blockchain submission failed — event stored in postgres only",
			"audit_event_id", event.ID, "err", err)
		return nil
	}

	// Step 3: Write tx_hash back to the postgres record
	if err := l.updateTxHash(ctx, event.ID, txHash.Hex()); err != nil {
		slog.Error("audit: failed to persist tx_hash", "err", err)
	}
	event.TxHash = txHash.Hex()

	slog.Info("audit: event logged to blockchain",
		"audit_event_id", event.ID,
		"tx_hash", txHash.Hex())

	return nil
}

// logToBlockchain submits the audit event to the smart contract
func (l *BlockchainAuditLogger) logToBlockchain(ctx context.Context, event *entity.AuditEvent) (*common.Hash, error) {
	// Convert item ID to *big.Int
	itemID := new(big.Int).SetUint64(event.ItemID)

	// Call smart contract's LogEvent method
	tx, err := l.contract.LogEvent(l.transactOpts, itemID, string(event.Action), event.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, l.ethClient, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for mining: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	txHash := receipt.TxHash
	return &txHash, nil
}

// updateTxHash updates the transaction hash for an audit event in PostgreSQL
func (l *BlockchainAuditLogger) updateTxHash(ctx context.Context, eventID uint64, txHash string) error {
	if pgLogger, ok := l.db.(*postgres.PostgresAuditLogger); ok {
		query := `UPDATE audit_events SET tx_hash = $1 WHERE id = $2`
		_, err := pgLogger.GetPool().Exec(ctx, query, txHash, int64(eventID))
		return err
	}
	return nil
}

// ListByItem delegates to the PostgreSQL implementation.
func (l *BlockchainAuditLogger) ListByItem(
	ctx context.Context,
	itemID uint64,
) ([]*entity.AuditEvent, error) {
	return l.db.ListByItem(ctx, itemID)
}

// GetContractAddress returns the deployed contract address
func (l *BlockchainAuditLogger) GetContractAddress() common.Address {
	return l.contractAddress
}
