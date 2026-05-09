// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// ItemStatus represents the current lifecycle state of an inventory item.
type ItemStatus string

const (
	// StatusActive marks an item currently in use or in storage.
	StatusActive ItemStatus = "active"
	// StatusInRepair marks an item temporarily out of service for maintenance.
	StatusInRepair ItemStatus = "in_repair"
	// StatusPendingDisposal marks an item awaiting disposal document upload.
	StatusPendingDisposal ItemStatus = "pending_disposal"
	// StatusDisposed marks an item that has been permanently written off.
	StatusDisposed ItemStatus = "disposed"
)

// ItemRelation represents a link between two items (e.g., computer + monitor).
type ItemRelation struct {
	ID            uint64    `json:"id"`
	ItemID        uint64    `json:"item_id"`
	RelatedItemID uint64    `json:"related_item_id"`
	RelatedItem   *Item     `json:"related_item,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     uint64    `json:"created_by"`
	Creator       *User     `json:"creator,omitempty"`
}

// Item is the central inventory record tracked by the system.
type Item struct {
	ID uint64 `json:"id"`

	Barcode         string `json:"barcode"`
	InventoryNumber string `json:"inventory_number"`
	Name            string `json:"name"`

	CategoryID uint64    `json:"category_id"`
	Category   *Category `json:"category,omitempty"`

	RoomID uint64 `json:"room_id"`
	Room   *Room  `json:"room,omitempty"`

	Description string `json:"description"`

	Status ItemStatus `json:"status"`

	TxHash string `json:"tx_hash,omitempty"`

	Photos []*ItemPhoto `json:"photos,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PendingDisposalAt *time.Time `json:"pending_disposal_at,omitempty"`
	DisposedAt        *time.Time `json:"disposed_at,omitempty"`

	CreatedBy    uint64 `json:"created_by"`
	LastEditedBy uint64 `json:"last_edited_by"`

	Creator    *User `json:"creator,omitempty"`
	LastEditor *User `json:"last_editor,omitempty"`
}

// IsActive returns true when the item is active.
func (item *Item) IsActive() bool {
	return item.Status == StatusActive
}

// IsMutable returns true when the item can be edited.
// Active and in_repair items are mutable — pending_disposal and disposed items are locked.
func (item *Item) IsMutable() bool {
	return item.Status != StatusDisposed && item.Status != StatusPendingDisposal
}

// MarkPendingDisposal transitions the item to pending_disposal state.
func (item *Item) MarkPendingDisposal(actorID uint64) {
	now := time.Now()
	item.Status = StatusPendingDisposal
	item.PendingDisposalAt = &now
	item.LastEditedBy = actorID
	item.UpdatedAt = now
}

// FinalizeDisposal transitions the item from pending_disposal to disposed state.
func (item *Item) FinalizeDisposal(actorID uint64) {
	now := time.Now()
	item.Status = StatusDisposed
	item.DisposedAt = &now
	item.LastEditedBy = actorID
	item.UpdatedAt = now
}
