// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// ItemStatus represents the current lifecycle state of an inventory item.
type ItemStatus string

const (
	// StatusActive marks an item currently in use or in storage.
	StatusActive ItemStatus = "active"
	// StatusDisposed marks an item that has been written off.
	StatusDisposed ItemStatus = "disposed"
)

// Item is the central inventory record tracked by the system.
// Immutable fields (Barcode, Name, CategoryID, RoomID) are set on creation
// and protected by the use-case layer — only Description and PhotoURL
// may be changed after creation.
type Item struct {
	ID uint64 `json:"id"`

	// Barcode is the physical barcode label on the asset.
	Barcode string `json:"barcode"`

	// Name is the human-readable asset name (e.g. "Dell Monitor U2422H").
	Name string `json:"name"`

	CategoryID uint64    `json:"category_id"`
	Category   *Category `json:"category,omitempty"`

	RoomID uint64 `json:"room_id"`
	Room   *Room  `json:"room,omitempty"`

	// Description is a mutable free-text note about the asset.
	Description string `json:"description"`

	// PhotoURL points to the stored photo of the asset; mutable.
	PhotoURL string `json:"photo_url,omitempty"`

	Status ItemStatus `json:"status"`

	// TxHash holds the latest blockchain transaction hash for this item.
	// Populated after blockchain integration; empty until then.
	TxHash string `json:"tx_hash,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedBy is the ID of the user who registered the item.
	CreatedBy uint64 `json:"created_by"`
	// LastEditedBy is the ID of the user who last modified the item.
	LastEditedBy uint64 `json:"last_edited_by"`

	Creator    *User `json:"creator,omitempty"`
	LastEditor *User `json:"last_editor,omitempty"`
}

// IsActive returns true when the item has not been disposed.
func (item *Item) IsActive() bool {
	return item.Status == StatusActive
}

// IsMutable returns true when mutable fields may be updated.
func (item *Item) IsMutable() bool {
	return item.Status != StatusDisposed
}

// Dispose transitions the item to the disposed state.
func (item *Item) Dispose(actorID uint64) {
	item.Status = StatusDisposed
	item.LastEditedBy = actorID
	item.UpdatedAt = time.Now()
}
