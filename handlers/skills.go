package handlers

import (
	"encoding/json"
	"github.com/blacktalenthubs/go-service-api/database"
	"github.com/blacktalenthubs/go-service-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// SkillHandler manages HTTP requests for skill resources
type SkillHandler struct {
	db *database.PostgresDB
}

// NewSkillHandler creates a new skill handler
func NewSkillHandler(db *database.PostgresDB) *SkillHandler {
	return &SkillHandler{
		db: db,
	}
}

// GetAll returns all skills
func (h *SkillHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	skills, err := h.db.GetAllSkills()
	if err != nil {
		http.Error(w, "Failed to get skills: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// Get returns a specific skill by ID
func (h *SkillHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	skill, err := h.db.GetSkill(id)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "skill with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get skill: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skill)
}

// Create adds a new skill
func (h *SkillHandler) Create(w http.ResponseWriter, r *http.Request) {
	var skill models.Skill

	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if skill.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	createdSkill, err := h.db.CreateSkill(skill)
	if err != nil {
		http.Error(w, "Failed to create skill: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSkill)
}

// Update modifies an existing skill
func (h *SkillHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	var skill models.Skill
	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if skill.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	updatedSkill, err := h.db.UpdateSkill(id, skill)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "skill with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update skill: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSkill)
}

// Delete removes a skill
func (h *SkillHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteSkill(id); err != nil {
		// Check if it's a not found error
		if err.Error() == "skill with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err.Error() == "cannot delete skill with id "+strconv.Itoa(id)+" because it is assigned to consultants" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Failed to delete skill: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
