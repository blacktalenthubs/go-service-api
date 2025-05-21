package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/blacktalenthubs/go-service-api/models"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"time"
)

// PostgresDB wraps the SQL DB connection
type PostgresDB struct {
	db *sql.DB
}

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// New creates a new database connection
func New(config Config) (*PostgresDB, error) {
	// Connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	// Add this line after forming the connStr
	log.Printf("Connection string: %s", connStr)
	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize the database
	if err := initDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

// Initialize the database schema
func initDatabase(db *sql.DB) error {
	// Create tables if they don't exist
	_, err := db.Exec(`
        -- Skills table
        CREATE TABLE IF NOT EXISTS skills (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL UNIQUE,
            description TEXT
        );

        -- Consultants table
        CREATE TABLE IF NOT EXISTS consultants (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE
        );

        -- ConsultantSkills junction table for many-to-many
        CREATE TABLE IF NOT EXISTS consultant_skills (
            consultant_id INTEGER REFERENCES consultants(id) ON DELETE CASCADE,
            skill_id INTEGER REFERENCES skills(id) ON DELETE CASCADE,
            PRIMARY KEY (consultant_id, skill_id)
        );
    `)

	return err
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	return db.db.Close()
}

// Consultant methods

// GetConsultant retrieves a consultant by ID
func (db *PostgresDB) GetConsultant(id int) (models.Consultant, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Consultant{}, err
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Get consultant
	var consultant models.Consultant
	err = tx.QueryRowContext(
		ctx,
		"SELECT id, name, email FROM consultants WHERE id = $1",
		id,
	).Scan(&consultant.ID, &consultant.Name, &consultant.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Consultant{}, fmt.Errorf("consultant with id %d not found", id)
		}
		return models.Consultant{}, err
	}

	// Get consultant skills
	rows, err := tx.QueryContext(
		ctx,
		"SELECT skill_id FROM consultant_skills WHERE consultant_id = $1",
		id,
	)
	if err != nil {
		return models.Consultant{}, err
	}
	defer rows.Close()

	// Collect skill IDs
	var skillIDs []int
	for rows.Next() {
		var skillID int
		if err := rows.Scan(&skillID); err != nil {
			return models.Consultant{}, err
		}
		skillIDs = append(skillIDs, skillID)
	}
	consultant.SkillIDs = skillIDs

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.Consultant{}, err
	}

	return consultant, nil
}

// GetAllConsultants returns all consultants
func (db *PostgresDB) GetAllConsultants() ([]models.Consultant, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query all consultants
	rows, err := db.db.QueryContext(ctx, "SELECT id, name, email FROM consultants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect consultants
	var consultants []models.Consultant
	for rows.Next() {
		var c models.Consultant
		if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
			return nil, err
		}
		consultants = append(consultants, c)
	}

	// Check for errors after scanning
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get skills for each consultant concurrently
	for i, consultant := range consultants {
		// This could be optimized with a single join query
		skillRows, err := db.db.QueryContext(
			ctx,
			"SELECT skill_id FROM consultant_skills WHERE consultant_id = $1",
			consultant.ID,
		)
		if err != nil {
			return nil, err
		}

		var skillIDs []int
		for skillRows.Next() {
			var skillID int
			if err := skillRows.Scan(&skillID); err != nil {
				skillRows.Close()
				return nil, err
			}
			skillIDs = append(skillIDs, skillID)
		}
		skillRows.Close()

		// Update skills in the consultant object
		consultants[i].SkillIDs = skillIDs
	}

	return consultants, nil
}

// CreateConsultant adds a new consultant
func (db *PostgresDB) CreateConsultant(consultant models.Consultant) (models.Consultant, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Consultant{}, err
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Insert consultant
	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO consultants (name, email) VALUES ($1, $2) RETURNING id",
		consultant.Name, consultant.Email,
	).Scan(&consultant.ID)

	if err != nil {
		return models.Consultant{}, err
	}

	// Insert consultant skills
	for _, skillID := range consultant.SkillIDs {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO consultant_skills (consultant_id, skill_id) VALUES ($1, $2)",
			consultant.ID, skillID,
		)
		if err != nil {
			return models.Consultant{}, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.Consultant{}, err
	}

	return consultant, nil
}

// UpdateConsultant updates an existing consultant
func (db *PostgresDB) UpdateConsultant(id int, consultant models.Consultant) (models.Consultant, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Begin a transaction
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Consultant{}, err
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Check if consultant exists
	var exists bool
	err = tx.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM consultants WHERE id = $1)",
		id,
	).Scan(&exists)

	if err != nil {
		return models.Consultant{}, err
	}

	if !exists {
		return models.Consultant{}, fmt.Errorf("consultant with id %d not found", id)
	}

	// Update consultant
	_, err = tx.ExecContext(
		ctx,
		"UPDATE consultants SET name = $1, email = $2 WHERE id = $3",
		consultant.Name, consultant.Email, id,
	)
	if err != nil {
		return models.Consultant{}, err
	}

	// Update consultant ID
	consultant.ID = id

	// Delete all consultant skills
	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM consultant_skills WHERE consultant_id = $1",
		id,
	)
	if err != nil {
		return models.Consultant{}, err
	}

	// Insert updated consultant skills
	for _, skillID := range consultant.SkillIDs {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO consultant_skills (consultant_id, skill_id) VALUES ($1, $2)",
			id, skillID,
		)
		if err != nil {
			return models.Consultant{}, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return models.Consultant{}, err
	}

	return consultant, nil
}

// DeleteConsultant removes a consultant
func (db *PostgresDB) DeleteConsultant(id int) error {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Delete consultant (cascade will handle consultant_skills)
	result, err := db.db.ExecContext(
		ctx,
		"DELETE FROM consultants WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	// Check if consultant existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("consultant with id %d not found", id)
	}

	return nil
}

// GetConsultantsBySkill returns all consultants with a specific skill
func (db *PostgresDB) GetConsultantsBySkill(skillID int) ([]models.Consultant, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query consultants with specific skill
	rows, err := db.db.QueryContext(
		ctx,
		`SELECT c.id, c.name, c.email 
         FROM consultants c
         JOIN consultant_skills cs ON c.id = cs.consultant_id
         WHERE cs.skill_id = $1`,
		skillID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect consultants
	var consultants []models.Consultant
	for rows.Next() {
		var c models.Consultant
		if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
			return nil, err
		}
		consultants = append(consultants, c)
	}

	// Check for errors after scanning
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get all skills for each consultant
	for i, consultant := range consultants {
		skillRows, err := db.db.QueryContext(
			ctx,
			"SELECT skill_id FROM consultant_skills WHERE consultant_id = $1",
			consultant.ID,
		)
		if err != nil {
			return nil, err
		}

		var skillIDs []int
		for skillRows.Next() {
			var id int
			if err := skillRows.Scan(&id); err != nil {
				skillRows.Close()
				return nil, err
			}
			skillIDs = append(skillIDs, id)
		}
		skillRows.Close()

		// Update skills in the consultant object
		consultants[i].SkillIDs = skillIDs
	}

	return consultants, nil
}

// Skill methods

// GetSkill retrieves a skill by ID
func (db *PostgresDB) GetSkill(id int) (models.Skill, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Get skill
	var skill models.Skill
	err := db.db.QueryRowContext(
		ctx,
		"SELECT id, name, description FROM skills WHERE id = $1",
		id,
	).Scan(&skill.ID, &skill.Name, &skill.Description)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Skill{}, fmt.Errorf("skill with id %d not found", id)
		}
		return models.Skill{}, err
	}

	return skill, nil
}

// GetAllSkills returns all skills
func (db *PostgresDB) GetAllSkills() ([]models.Skill, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query all skills
	rows, err := db.db.QueryContext(ctx, "SELECT id, name, description FROM skills")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect skills
	var skills []models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.Description); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	// Check for errors after scanning
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

// CreateSkill adds a new skill
func (db *PostgresDB) CreateSkill(skill models.Skill) (models.Skill, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert skill
	err := db.db.QueryRowContext(
		ctx,
		"INSERT INTO skills (name, description) VALUES ($1, $2) RETURNING id",
		skill.Name, skill.Description,
	).Scan(&skill.ID)

	if err != nil {
		return models.Skill{}, err
	}

	return skill, nil
}

// UpdateSkill updates an existing skill
func (db *PostgresDB) UpdateSkill(id int, skill models.Skill) (models.Skill, error) {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Update skill
	result, err := db.db.ExecContext(
		ctx,
		"UPDATE skills SET name = $1, description = $2 WHERE id = $3",
		skill.Name, skill.Description, id,
	)
	if err != nil {
		return models.Skill{}, err
	}

	// Check if skill existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Skill{}, err
	}

	if rowsAffected == 0 {
		return models.Skill{}, fmt.Errorf("skill with id %d not found", id)
	}

	// Update skill ID
	skill.ID = id

	return skill, nil
}

// DeleteSkill removes a skill
func (db *PostgresDB) DeleteSkill(id int) error {
	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if skill is being used by any consultant
	var inUse bool
	err := db.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM consultant_skills WHERE skill_id = $1)",
		id,
	).Scan(&inUse)

	if err != nil {
		return err
	}

	if inUse {
		return fmt.Errorf("cannot delete skill with id %d because it is assigned to consultants", id)
	}

	// Delete skill
	result, err := db.db.ExecContext(
		ctx,
		"DELETE FROM skills WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	// Check if skill existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("skill with id %d not found", id)
	}

	return nil
}
