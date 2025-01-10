package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time // The timestamp when record was created.
	UpdatedAt time.Time // The timestamp when the record was last updated.
	DeletedAt time.Time // Soft delete field with an index for querying.

	Public      bool   `json:"public"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Address     string `json:"address"`
	Deadline    string `json:"deadline"`
	Budget      int    `json:"budget"`
	Status      string `json:"status"`
}
