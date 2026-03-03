// Package entity defines the core domain models of the inventory system.
package entity

// Room represents a physical location inside a building where items are kept.
type Room struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	BuildingID uint64    `json:"building_id"`
	Building   *Building `json:"building,omitempty"`
}
