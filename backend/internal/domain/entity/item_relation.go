// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// ItemRelation represents a symmetric relationship between two items.
// Items can be linked together to form kits or bundles (e.g., computer + monitor).
// The relationship is symmetric: if item A is linked to item B, then B is linked to A.
type ItemRelation struct {
	ID        uint64    `json:"id"`
	ItemID1   uint64    `json:"item_id_1"`
	ItemID2   uint64    `json:"item_id_2"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint64    `json:"created_by"`
}

// LinkedItem represents an item that is linked to another item.
// This is used when fetching related items for display.
type LinkedItem struct {
	*Item
	RelationID uint64    `json:"relation_id"`
	LinkedAt   time.Time `json:"linked_at"`
}
