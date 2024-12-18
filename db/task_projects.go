package db

import "github.com/jmoiron/sqlx"

const M007_TaskProjectsSchema = `
CREATE TABLE IF NOT EXISTS taskProjects (
	taskId INTEGER NOT NULL,
	projectId INTEGER NOT NULL,
	FOREIGN KEY (taskId) REFERENCES tasks(id) ON DELETE CASCADE,
	FOREIGN KEY (projectId) REFERENCES projects(id) ON DELETE CASCADE,
	UNIQUE (taskId, projectId),
	PRIMARY KEY (taskId, projectId)
);
PRAGMA user_version = 6;
`

// TaskProjectLink is a struct that represents a link between a task and a project
type TaskProjectLink struct {
	TaskID    int64 `db:"taskId"`
	ProjectID int64 `db:"projectId"`
}

// ProjectLinkTaskTx will link a task to a project
func (s *Store) ProjectLinkTaskTx(tx *sqlx.Tx, projectId, taskId int64) error {
	var _, err = tx.Exec(`INSERT INTO taskProjects (projectId, taskId) VALUES (?, ?)`, projectId, taskId)
	return err
}

// ProjectUnlinkTaskTx will unlink a task from a project
func (s *Store) ProjectTasksList() ([]TaskProjectLink, error) {
	var links []TaskProjectLink
	err := s.Con.Select(&links, `SELECT * FROM taskProjects`)
	return links, err
}
