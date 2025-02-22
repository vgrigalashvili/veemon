package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id"` // The unique identifier of the task.
	CreatedAt time.Time // The timestamp when record was created.
	UpdatedAt time.Time // The timestamp when the record was last updated.
	DeletedAt time.Time // Soft delete field with an index for querying.

	Public      bool   `json:"public"`      // Indicates if the task is public or private.
	Title       string `json:"title"`       // Title of the task
	Description string `json:"description"` // Description of the task
	Location    string `json:"location"`    // Location of the task
	Address     string `json:"address"`     // Address of the task
	DeadLine    string `json:"dead_line"`   // Deadline of the task
	Budget      int    `json:"budget"`      // Budget of the task
	Status      string `json:"status"`      // Status of the task
}
