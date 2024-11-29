package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const M010_TaskDependenciesSchema = `
CREATE TABLE IF NOT EXISTS taskDependencies (
	taskId INTEGER NOT NULL,
	dependsOnId INTEGER NOT NULL,
	FOREIGN KEY (taskId) REFERENCES tasks(id) ON DELETE CASCADE,
	FOREIGN KEY (dependsOnId) REFERENCES tasks(id) ON DELETE CASCADE,
	UNIQUE (taskId, dependsOnId)
);
PRAGMA user_version = 10;
`

type TaskDependency struct {
	TaskID      int64 `db:"taskId"`
	DependsOnID int64 `db:"dependsOnId"`
}

func (store *Store) TaskDependsOnTx(tx *sqlx.Tx, taskId int64, dependsOnId int64) error {
	var _, err = tx.Exec(`INSERT INTO taskDependencies (taskId, dependsOnId) VALUES (?, ?)`, taskId, dependsOnId)
	if err != nil {
		return fmt.Errorf("Failed to insert task dependency: %w", err)
	}
	return nil
}

func (store *Store) GetDependenciesForTask(taskId int64) ([]TaskDependency, error) {
	var deps []TaskDependency
	err := store.Con.Select(&deps, `SELECT * FROM taskDependencies WHERE taskId = ?`, taskId)
	return deps, err
}

func (store *Store) DeleteDependenciesForCompletedTask(completedTaskId int64) error {
	_, err := store.Con.Exec(`DELETE FROM taskDependencies WHERE taskId = ? OR dependsOnId = ?`, completedTaskId, completedTaskId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil
		}
	}
	return err
}
