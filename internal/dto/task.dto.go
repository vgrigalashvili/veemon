package dto

import "github.com/vgrigalashvili/veemon/internal/model"

// CreateTask is a data transfer object for creating a new task.
type CreateTask struct {
	Private     bool               `json:"private" validate:"omitempty"`     // Private task or not, optional.
	Title       string             `json:"title" validate:"required,min=3"`  // Title of the task, required.
	Description string             `json:"description" validate:"omitempty"` // Description of the task, optional.
	Category    string             `json:"category" validate:"omitempty"`    // Category of the task, optional.
	Location    string             `json:"location" validate:"omitempty"`    // Location of the task, optional.
	Deadline    string             `json:"deadline" validate:"omitempty"`    // Deadline of the task, optional.
	Budget      int                `json:"budget" validate:"required"`       // Budget of the task, required.
	Attachments []model.Attachment `json:"attachments" validate:"omitempty"` // Attachments to the task, optional.
}
