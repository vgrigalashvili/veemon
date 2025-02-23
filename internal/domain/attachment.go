package domain

import (
	"time"

	"github.com/google/uuid"
)

// task attachment

type Attachment struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}
