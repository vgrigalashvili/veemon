package dto

import "github.com/vgrigalashvili/veemon/internal/domain"

type CreateTask struct {
	Private     bool                `json:"private" validate:"omitempty"`
	Title       string              `json:"title" validate:"required,min=3"`
	Description string              `json:"description" validate:"omitempty"`
	Category    string              `json:"category" validate:"omitempty"`
	Location    string              `json:"location" validate:"omitempty"`
	Deadline    string              `json:"deadline" validate:"omitempty"`
	Budget      int                 `json:"budget" validate:"required"`
	Attachments []domain.Attachment `json:"attachments" validate:"omitempty"`
}
