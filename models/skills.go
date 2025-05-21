package models

// Skill represents a skill that consultants can have
type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
