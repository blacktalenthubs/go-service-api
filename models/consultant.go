package models

// Consultant represents a consultant in the system
type Consultant struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SkillIDs  []int  `json:"skill_ids"`
	ProjectID *int   `json:"project_id,omitempty"`
}

// ConsultantSkill represents the many-to-many relationship
type ConsultantSkill struct {
	ConsultantID int `json:"consultant_id"`
	SkillID      int `json:"skill_id"`
}
