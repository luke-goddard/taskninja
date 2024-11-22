package db

import "database/sql"

const M000_TaskSchema = `
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT,
	completed INTEGER NOT NULL DEFAULT 0,
	priority INTEGER NOT NULL DEFAULT 0 CHECK (priority >= 0 AND priority <= 3),
	dueUtc TEXT,
	updatedAtUtc TEXT NOT NULL DEFAULT current_timestamp,
	createdAtUtc TEXT NOT NULL DEFAULT current_timestamp,
	completedAtUtc TEXT
	startedAtUtc TEXT
);
PRAGMA user_version = 0;
`

const M003_TaskSchema = `
ALTER TABLE tasks ADD COLUMN startedAtUtc TEXT;
PRAGMA user_version = 3;
`

type TaskPriority int

const (
	TaskPriorityNone TaskPriority = iota
	TaskPriorityLow
	TaskPriorityMedium
	TaskPriorityHigh
)

type Task struct {
	ID           int           `json:"id" db:"id"`
	Title        string        `json:"title" db:"title"`
	Description  *string       `json:"description" db:"description"`
	Due          *string       `json:"due" db:"dueUtc"`
	Completed    *bool         `json:"completed" db:"completed"`
	Priority     *TaskPriority `json:"priority" db:"priority"`
	CreatedUtc   *string       `json:"createdUtc" db:"createdAtUtc"`
	UpdatedAtUtc *string       `json:"updatedAtUtc" db:"updatedAtUtc"`
	CompletedUtc *string       `json:"completedUtc" db:"completedAtUtc"`
	StartedUtc   *string       `json:"startedUtc" db:"startedAtUtc"`
}

func (store *Store) ListTasks() ([]Task, error) {
	var tasks []Task
	err := store.Con.Select(&tasks, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (store *Store) DeleteTaskById(id int64) (bool, error) {
	var err error
	var res sql.Result
	var rowsAffected int64
	res, err = store.Con.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func (store *Store) StartTaskById(id int64) (bool, error) {
	var err error
	var res sql.Result
	var rowsAffected int64
	res, err = store.Con.Exec("UPDATE tasks SET completed = 0 WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}
