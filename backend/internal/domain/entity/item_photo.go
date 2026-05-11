// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// ItemPhoto represents a single photo attached to an inventory item.
type ItemPhoto struct {
	ID         uint64    `json:"id"`
	ItemID     uint64    `json:"item_id"`
	Base64Data string    `json:"base64_data"`
	MimeType   string    `json:"mime_type"`
	CreatedAt  time.Time `json:"created_at"`
}
