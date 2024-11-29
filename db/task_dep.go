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

func (store *Store) TaskDependsOnTx(tx *sqlx.Tx, taskId int64, dependsOnId int64) (error) {
	var _, err = tx.Exec(`INSERT INTO taskDependencies (taskId, dependsOnId) VALUES (?, ?)`, taskId, dependsOnId)
	if err != nil {
		return  fmt.Errorf("Failed to insert task dependency: %w", err)
	}
	return nil
}
