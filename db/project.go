package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const M006_ProjectSchema = `
CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL
);
PRAGMA user_version = 6;
`

// A project may be assigned to a task, and that project may be multiple words.
type Project struct {
	ID    int64  `db:"id"`    // Unique identifier
	Title string `db:"title"` // Project title
}

// ProjectGetIDByNameOrCreate will get the project ID by name or create it if it does not exist.
func (s *Store) ProjectGetIDByNameOrCreateTx(tx *sqlx.Tx, title string) (int64, error) {
	var id int64
	var err = tx.Get(&id, `SELECT id FROM projects WHERE title = ?`, title)
	if errors.Is(err, sql.ErrNoRows) {
		result, err := tx.Exec(`INSERT INTO projects (title) VALUES (?)`, title)
		if err != nil {
			return 0, fmt.Errorf("Failed to insert project: %w", err)
		}
		return result.LastInsertId()
	}
	if err != nil {
		return 0, fmt.Errorf("Failed to get project: %w", err)
	}
	return id, err
}

// ListProjects returns a list of all projects.
func (s *Store) ListProjects() ([]Project, error) {
	var projects []Project
	err := s.Con.Select(&projects, `SELECT * FROM projects`)
	return projects, err
}
