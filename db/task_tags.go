package db

const M002_TaskTagsSchema = `
CREATE TABLE IF NOT EXISTS taskTags (
	taskID INTEGER NOT NULL,
	tagID INTEGER NOT NULL,
	PRIMARY KEY (taskID, tagID)
);`

type TaskTag struct {
	TaskID int `json:"taskID" db:"taskID"`
	TagID  int `json:"tagID" db:"tagID"`
}
