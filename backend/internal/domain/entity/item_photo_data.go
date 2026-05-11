// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// ItemPhoto represents a single photo attached to an inventory item.
type ItemPhoto struct {
	ID        uint64    `json:"id"`
	ItemID    uint64    `json:"item_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// ItemPhotoData represents the Base64-encoded binary data for a photo.
// This is stored in a separate table to keep main queries lightweight.
type ItemPhotoData struct {
	ID        uint64    `json:"id"`
	PhotoID   uint64    `json:"photo_id"`
	Data      []byte    `json:"data"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}
