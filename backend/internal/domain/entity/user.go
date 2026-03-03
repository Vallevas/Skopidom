// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// UserRole represents the access level granted to a user.
type UserRole string

const (
	// RoleAdmin has full access: create/edit/delete items and manage users.
	RoleAdmin UserRole = "admin"
	// RoleEditor can create and edit items but cannot manage users.
	RoleEditor UserRole = "editor"
)

// User represents a system account that can interact with inventory records.
type User struct {
	ID           uint64    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsAdmin returns true when the user holds the admin role.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanEdit returns true when the user is permitted to modify inventory records.
func (u *User) CanEdit() bool {
	return u.Role == RoleAdmin || u.Role == RoleEditor
}
