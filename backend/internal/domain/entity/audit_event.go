// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// AuditAction represents the type of change made to an inventory item.
type AuditAction string

const (
	// ActionCreated is recorded when a new item is registered.
	ActionCreated AuditAction = "created"
	// ActionUpdated is recorded when an item's mutable fields are changed.
	ActionUpdated AuditAction = "updated"
	// ActionDisposed is recorded when an item is permanently written off.
	ActionDisposed AuditAction = "disposed"
	// ActionMoved is recorded when an item is moved to a different room.
	// Payload contains from_room and to_room fields for full traceability.
	ActionMoved AuditAction = "moved"
	// ActionSentToRepair is recorded when an item's status changes to in_repair.
	ActionSentToRepair AuditAction = "sent_to_repair"
	// ActionReturnedFromRepair is recorded when an in_repair item returns to active.
	ActionReturnedFromRepair AuditAction = "returned_from_repair"
	// ActionPendingDisposal is recorded when an item enters pending_disposal status.
	ActionPendingDisposal AuditAction = "pending_disposal"
	// ActionDisposalFinalized is recorded when disposal is finalized with documents.
	ActionDisposalFinalized AuditAction = "disposal_finalized"
)

// AuditCategory represents the logical grouping of audit actions.
// Actions are split into two categories for UI presentation:
// - StatusLog: lifecycle events (creation, movement, disposal)
// - Changelog: modifications and repairs
type AuditCategory string

const (
	// CategoryStatusLog groups actions related to item lifecycle and location changes.
	CategoryStatusLog AuditCategory = "status_log"
	// CategoryChangelog groups actions related to item modifications and repairs.
	CategoryChangelog AuditCategory = "changelog"
)

// Category returns the logical category of this audit action.
// This categorization is used for splitting the audit log into two separate views:
// - Status Log: tracks lifecycle events (creation, movement, disposal)
// - Changelog: tracks modifications and repairs
func (a AuditAction) Category() AuditCategory {
	switch a {
	case ActionCreated, ActionMoved, ActionDisposed, ActionPendingDisposal, ActionDisposalFinalized:
		return CategoryStatusLog
	case ActionUpdated, ActionSentToRepair, ActionReturnedFromRepair:
		return CategoryChangelog
	default:
		// Default to changelog for any unknown actions
		return CategoryChangelog
	}
}

// AuditEvent is an immutable record of a single lifecycle change on an item.
type AuditEvent struct {
	ID      uint64 `json:"id"`
	ItemID  uint64 `json:"item_id"`
	Item    *Item  `json:"-"`
	ActorID uint64 `json:"actor_id"`
	Actor   *User  `json:"actor,omitempty"`

	Action AuditAction `json:"action"`

	// Payload is a JSON snapshot of the item state at the moment of the event.
	// For ActionMoved it additionally contains from/to room and building info.
	Payload string `json:"payload"`

	// TxHash is the Ethereum transaction hash anchoring this event on-chain.
	// Empty when blockchain integration is not enabled.
	TxHash string `json:"tx_hash,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
