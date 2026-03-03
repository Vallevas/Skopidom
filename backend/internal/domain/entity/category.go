// Package entity defines the core domain models of the inventory system.
package entity

// Category classifies inventory items (e.g. "Monitor", "PC", "Furniture").
type Category struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
