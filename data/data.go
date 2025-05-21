package data

import (
	"fmt"
	"github.com/blacktalenthubs/go-service-api/models"
	"sync"
)

// Store provides an in-memory data store with thread-safe operations
type Store struct {
	consultants map[int]models.Consultant
	skills      map[int]models.Skill
	projects    map[int]models.Project
	mutex       sync.RWMutex

	// Auto-incrementing IDs
	nextConsultantID int
	nextSkillID      int
	nextProjectID    int
}

// NewStore creates and initializes a new data store
func NewStore() *Store {
	store := &Store{
		consultants:      make(map[int]models.Consultant),
		skills:           make(map[int]models.Skill),
		projects:         make(map[int]models.Project),
		nextConsultantID: 1,
		nextSkillID:      1,
		nextProjectID:    1,
	}

	// Initialize with sample data
	store.seedData()
	return store
}

func (s *Store) seedData() {
	// Add skills
	programming := s.CreateSkill(models.Skill{Name: "Programming", Description: "Software development skills"})
	projectManagement := s.CreateSkill(models.Skill{Name: "Project Management", Description: "Managing project timelines and resources"})
	dataAnalysis := s.CreateSkill(models.Skill{Name: "Data Analysis", Description: "Analyzing and interpreting complex data"})

	// Add projects
	webApp := s.CreateProject(models.Project{Name: "Web Application", Description: "Customer portal application", ClientName: "Acme Inc"})
	dataWarehouse := s.CreateProject(models.Project{Name: "Data Warehouse", Description: "Data warehouse implementation", ClientName: "BigData Corp"})

	// Add consultants
	s.CreateConsultant(models.Consultant{Name: "John Doe", Email: "john@example.com", SkillIDs: []int{programming.ID, projectManagement.ID}, ProjectID: &webApp.ID})
	s.CreateConsultant(models.Consultant{Name: "Jane Smith", Email: "jane@example.com", SkillIDs: []int{dataAnalysis.ID}, ProjectID: &dataWarehouse.ID})
	s.CreateConsultant(models.Consultant{Name: "Bob Johnson", Email: "bob@example.com", SkillIDs: []int{programming.ID, dataAnalysis.ID}})
}

// Consultant operations

// GetConsultant retrieves a consultant by ID
func (s *Store) GetConsultant(id int) (models.Consultant, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	consultant, exists := s.consultants[id]
	if !exists {
		return models.Consultant{}, fmt.Errorf("consultant with id %d not found", id)
	}

	return consultant, nil
}

// GetAllConsultants returns all consultants
func (s *Store) GetAllConsultants() []models.Consultant {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	consultants := make([]models.Consultant, 0, len(s.consultants))
	for _, consultant := range s.consultants {
		consultants = append(consultants, consultant)
	}

	return consultants
}

func (s *Store) CreateConsultant(consultant models.Consultant) models.Consultant {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Assign ID
	consultant.ID = s.nextConsultantID
	s.nextConsultantID++

	// Store consultant
	s.consultants[consultant.ID] = consultant

	return consultant
}

func (s *Store) UpdateConsultant(id int, consultant models.Consultant) (models.Consultant, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.consultants[id]; !exists {
		return models.Consultant{}, fmt.Errorf("consultant with id %d not found", id)
	}

	// Ensure ID doesn't change
	consultant.ID = id

	// Update consultant
	s.consultants[id] = consultant

	return consultant, nil
}

func (s *Store) DeleteConsultant(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.consultants[id]; !exists {
		return fmt.Errorf("consultant with id %d not found", id)
	}

	delete(s.consultants, id)
	return nil
}

func (s *Store) GetConsultantsBySkill(skillID int) []models.Consultant {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var consultants []models.Consultant

	for _, consultant := range s.consultants {
		for _, id := range consultant.SkillIDs {
			if id == skillID {
				consultants = append(consultants, consultant)
				break
			}
		}
	}

	return consultants
}

// Skill operations

// GetSkill retrieves a skill by ID
func (s *Store) GetSkill(id int) (models.Skill, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	skill, exists := s.skills[id]
	if !exists {
		return models.Skill{}, fmt.Errorf("skill with id %d not found", id)
	}

	return skill, nil
}

// GetAllSkills returns all skills
func (s *Store) GetAllSkills() []models.Skill {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	skills := make([]models.Skill, 0, len(s.skills))
	for _, skill := range s.skills {
		skills = append(skills, skill)
	}

	return skills
}

// CreateSkill adds a new skill
func (s *Store) CreateSkill(skill models.Skill) models.Skill {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Assign ID
	skill.ID = s.nextSkillID
	s.nextSkillID++

	// Store skill
	s.skills[skill.ID] = skill

	return skill
}

// UpdateSkill updates an existing skill
func (s *Store) UpdateSkill(id int, skill models.Skill) (models.Skill, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.skills[id]; !exists {
		return models.Skill{}, fmt.Errorf("skill with id %d not found", id)
	}

	// Ensure ID doesn't change
	skill.ID = id

	// Update skill
	s.skills[id] = skill

	return skill, nil
}

// DeleteSkill removes a skill
func (s *Store) DeleteSkill(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.skills[id]; !exists {
		return fmt.Errorf("skill with id %d not found", id)
	}

	delete(s.skills, id)
	return nil
}

// Project operations

// GetProject retrieves a project by ID
func (s *Store) GetProject(id int) (models.Project, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	project, exists := s.projects[id]
	if !exists {
		return models.Project{}, fmt.Errorf("project with id %d not found", id)
	}

	return project, nil
}

// GetAllProjects returns all projects
func (s *Store) GetAllProjects() []models.Project {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	projects := make([]models.Project, 0, len(s.projects))
	for _, project := range s.projects {
		projects = append(projects, project)
	}

	return projects
}

// CreateProject adds a new project
func (s *Store) CreateProject(project models.Project) models.Project {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Assign ID
	project.ID = s.nextProjectID
	s.nextProjectID++

	// Store project
	s.projects[project.ID] = project

	return project
}

// UpdateProject updates an existing project
func (s *Store) UpdateProject(id int, project models.Project) (models.Project, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.projects[id]; !exists {
		return models.Project{}, fmt.Errorf("project with id %d not found", id)
	}

	// Ensure ID doesn't change
	project.ID = id

	// Update project
	s.projects[id] = project

	return project, nil
}

// DeleteProject removes a project
func (s *Store) DeleteProject(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.projects[id]; !exists {
		return fmt.Errorf("project with id %d not found", id)
	}

	delete(s.projects, id)

	// Update any consultants assigned to this project
	for consultantID, consultant := range s.consultants {
		if consultant.ProjectID != nil && *consultant.ProjectID == id {
			consultant.ProjectID = nil
			s.consultants[consultantID] = consultant
		}
	}

	return nil
}
