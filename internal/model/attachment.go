package model

import (
	"time"

	"github.com/google/uuid"
)

// task attachment

type Attachment struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time // The timestamp when the record was created.
	UpdatedAt time.Time // The timestamp when the record was last updated.
	DeletedAt time.Time // Soft delete field with an index for querying.

	Name string `json:"name"` // The name of the attachment.
	Type string `json:"type"` // The type of the attachment.
	Size int    `json:"size"` // The size of the attachment.
	// OwnerID   uuid.UUID `json:"owner_id"`  // The owner of the attachment.

}
