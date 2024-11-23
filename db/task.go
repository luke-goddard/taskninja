package db

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

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
	ID           int64          `json:"id" db:"id"`
	Title        string         `json:"title" db:"title"`
	Description  *string        `json:"description" db:"description"`
	Due          *string        `json:"due" db:"dueUtc"`
	Completed    bool           `json:"completed" db:"completed"`
	Priority     TaskPriority   `json:"priority" db:"priority"`
	CreatedUtc   *string        `json:"createdUtc" db:"createdAtUtc"`
	UpdatedAtUtc *string        `json:"updatedAtUtc" db:"updatedAtUtc"`
	CompletedUtc *string        `json:"completedUtc" db:"completedAtUtc"`
	StartedUtc   sql.NullString `json:"startedUtc" db:"startedAtUtc"`
}

func (task *Task) IsStarted() bool {
	return task.StartedUtc.Valid
}

func (task *Task) TimeSinceStarted() time.Duration {
	if !task.IsStarted() {
		return 0
	}
	var startedAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", task.StartedUtc.String)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse startedAt")
		return 0
	}
	return time.Since(startedAt)
}

func (task *Task) TimeSinceStartedStr() string {
	return task.TimeSinceStarted().String()
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

func (store *Store) StartTaskById(id int64) (*Task, error) {
	var sql = `
	UPDATE tasks
	SET
		updatedAtUtc = current_timestamp,
		startedAtUtc = case when startedAtUtc is null then ? else startedAtUtc end
	WHERE id = ?
	RETURNING *
	`
	var task = &Task{}
	var now = time.Now().UTC().String()
	var row = store.Con.QueryRowx(sql, now, id)
	var err = row.StructScan(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (store *Store) CreateTask(task *Task) (*Task, error) {
	var sql = `
	INSERT INTO tasks
		(
			title, description, dueUtc,
			completed, priority, createdAtUtc,
			updatedAtUtc, completedAtUtc, startedAtUtc
		)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING *
	`
	var newTask = &Task{}
	var row = store.Con.QueryRowx(
		sql,
		task.Title, task.Description, task.Due,
		task.Completed, task.Priority, time.Now().UTC().String(),
		time.Now().UTC().String(), task.CompletedUtc, task.StartedUtc,
	)
	var err = row.StructScan(newTask)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}
