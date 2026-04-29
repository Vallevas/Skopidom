// Package entity defines the core domain models of the inventory system.
package entity

import "time"

// DisposalDocument represents a document attached to an item during disposal process.
type DisposalDocument struct {
	ID         uint64    `json:"id"`
	ItemID     uint64    `json:"item_id"`
	Filename   string    `json:"filename"`
	URL        string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
	UploadedBy uint64    `json:"uploaded_by"`
	Uploader   *User     `json:"uploader,omitempty"`
}
