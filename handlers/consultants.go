package handlers

import (
	"encoding/json"
	"github.com/blacktalenthubs/go-service-api/database"
	"github.com/blacktalenthubs/go-service-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// ConsultantHandler manages HTTP requests for consultant resources
type ConsultantHandler struct {
	db *database.PostgresDB
}

// NewConsultantHandler creates a new consultant handler
func NewConsultantHandler(db *database.PostgresDB) *ConsultantHandler {
	return &ConsultantHandler{
		db: db,
	}
}

// GetAll returns all consultants
func (h *ConsultantHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	consultants, err := h.db.GetAllConsultants()
	if err != nil {
		http.Error(w, "Failed to get consultants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultants)
}

// Get returns a specific consultant by ID
func (h *ConsultantHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	consultant, err := h.db.GetConsultant(id)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "consultant with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get consultant: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultant)
}

// Create adds a new consultant
func (h *ConsultantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var consultant models.Consultant

	if err := json.NewDecoder(r.Body).Decode(&consultant); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if consultant.Name == "" || consultant.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	createdConsultant, err := h.db.CreateConsultant(consultant)
	if err != nil {
		http.Error(w, "Failed to create consultant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdConsultant)
}

// Update modifies an existing consultant
func (h *ConsultantHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	var consultant models.Consultant
	if err := json.NewDecoder(r.Body).Decode(&consultant); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if consultant.Name == "" || consultant.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	updatedConsultant, err := h.db.UpdateConsultant(id, consultant)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "consultant with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update consultant: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedConsultant)
}

// Delete removes a consultant
func (h *ConsultantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid consultant ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteConsultant(id); err != nil {
		// Check if it's a not found error
		if err.Error() == "consultant with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete consultant: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetBySkill returns all consultants with a specific skill
func (h *ConsultantHandler) GetBySkill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skillID, err := strconv.Atoi(vars["skill_id"])
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	consultants, err := h.db.GetConsultantsBySkill(skillID)
	if err != nil {
		http.Error(w, "Failed to get consultants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consultants)
}
