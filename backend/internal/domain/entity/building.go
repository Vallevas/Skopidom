// Package entity defines the core domain models of the inventory system.
package entity

// Building represents a physical university building that contains rooms.
type Building struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
