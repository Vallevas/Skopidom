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
	// ActionDisposed is recorded when an item is written off.
	ActionDisposed AuditAction = "disposed"
	// ActionMoved is recorded when an item is moved to a different room.
	// The Payload contains from_room and to_room fields for full traceability.
	ActionMoved AuditAction = "moved"
)

// AuditEvent is an immutable record of a single lifecycle change on an item.
type AuditEvent struct {
	ID      uint64 `json:"id"`
	ItemID  uint64 `json:"item_id"`
	Item    *Item  `json:"-"`
	ActorID uint64 `json:"actor_id"`
	Actor   *User  `json:"actor,omitempty"`

	Action AuditAction `json:"action"`

	// Payload is a JSON snapshot of the item state at the moment of the event.
	// For ActionMoved it additionally contains from_room_id, from_room_name,
	// from_building_name, to_room_id, to_room_name, to_building_name.
	Payload string `json:"payload"`

	// TxHash is the Ethereum transaction hash anchoring this event on-chain.
	// Empty when blockchain integration is not enabled.
	TxHash string `json:"tx_hash,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
