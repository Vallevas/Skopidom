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
	// StatusDisposed marks an item that has been permanently written off.
	StatusDisposed ItemStatus = "disposed"
)

// Item is the central inventory record tracked by the system.
type Item struct {
	ID uint64 `json:"id"`

	Barcode string `json:"barcode"`
	Name    string `json:"name"`

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
// Both active and in_repair items are mutable — only disposed items are locked.
func (item *Item) IsMutable() bool {
	return item.Status != StatusDisposed
}

// Dispose transitions the item to the disposed state.
func (item *Item) Dispose(actorID uint64) {
	item.Status = StatusDisposed
	item.LastEditedBy = actorID
	item.UpdatedAt = time.Now()
}
