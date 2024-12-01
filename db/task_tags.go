package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const M002_TaskTagsSchema = `
CREATE TABLE IF NOT EXISTS taskTags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	UNIQUE (taskID, tagID),
	PRIMARY KEY(taskID, tagID),
	FOREIGN KEY(taskID) REFERENCES tasks(id) ON DELETE CASCADE,
	FOREIGN KEY(tagID) REFERENCES tags(id) ON DELETE CASCADE
);
PRAGMA user_version = 2;
`

// Used to link tasks and tags together
type TaskTag struct {
	ID     int `json:"id" db:"id"`         // Unique identifier of the task-tag link
	TaskID int `json:"taskID" db:"taskID"` // ID of the task
	TagID  int `json:"tagID" db:"tagID"`   // ID of the tag
}

// TagLinkTask will link a task to a tag
func (store *Store) TagLinkTaskCtx(ctx context.Context, taskId, tagId int64) error {
	_, err := store.Con.ExecContext(ctx, "INSERT INTO taskTags (taskID, tagID) VALUES (?, ?)", taskId, tagId)
	if err != nil {
		return fmt.Errorf("Failed to link task and tag: %w", err)
	}
	return nil
}

// TagLinkTaskTx will link a task to a tag inside of a transaction
func (store *Store) TagLinkTaskTx(tx *sqlx.Tx, taskId, tagId int64) error {
	_, err := tx.Exec("INSERT INTO taskTags (taskID, tagID) VALUES (?, ?)", taskId, tagId)
	if err != nil {
		return fmt.Errorf("Failed to link task and tag: %w", err)
	}
	return nil
}

// TagUnlinkTask will unlink a tag
func (store *Store) TagUnlinkTask(taskId, tagId int64) error {
	_, err := store.Con.Exec("DELETE FROM taskTags WHERE taskID = ? AND tagID = ?", taskId, tagId)
	if err != nil {
		return fmt.Errorf("Failed to unlink task and tag: %w", err)
	}
	return nil
}

// TagUnlinkTask will unlink a tag inside of a transaction
func (store *Store) TagUnlinkTaskTx(tx *sqlx.Tx, taskId, tagId int64) error {
	_, err := tx.Exec("DELETE FROM taskTags WHERE taskID = ? AND tagID = ?", taskId, tagId)
	if err != nil {
		return fmt.Errorf("Failed to unlink task and tag: %w", err)
	}
	return nil
}
