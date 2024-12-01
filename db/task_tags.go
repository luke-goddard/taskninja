package db

import "fmt"

const M002_TaskTagsSchema = `
CREATE TABLE IF NOT EXISTS taskTags (
	taskID INTEGER NOT NULL,
	tagID INTEGER NOT NULL,
	UNIQUE (taskID, tagID),
	PRIMARY KEY (taskID, tagID)
);
-- PRAGMA user_version = 2;
`

type TaskTag struct {
	TaskID int `json:"taskID" db:"taskID"`
	TagID  int `json:"tagID" db:"tagID"`
}

func (store *Store) TagLinkTask(taskId, tagId int64) error {
	_, err := store.Con.Exec("INSERT INTO taskTags (taskID, tagID) VALUES (?, ?)", taskId, tagId)
	if err != nil {
		return fmt.Errorf("Failed to link task and tag: %w", err)
	}
	return nil
}
