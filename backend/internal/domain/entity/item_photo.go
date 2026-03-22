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
