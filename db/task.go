package db

import (
	"database/sql"
	"strings"
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
);
PRAGMA user_version = 0;
`

const M003_TaskSchema = `
ALTER TABLE tasks ADD COLUMN startedAtUtc TEXT;
PRAGMA user_version = 3;
`

const M004_TaskSchema = `
ALTER TABLE tasks ADD COLUMN state INTEGER NOT NULL DEFAULT 0 CHECK (state >= 0 AND state <= 2);
PRAGMA user_version = 4;
`

const M005_TaskSchema = `
ALTER TABLE tasks DROP COLUMN completed;
PRAGMA user_version = 5;
`

type TaskPriority int

const (
	TaskPriorityNone TaskPriority = iota // Default priority
	TaskPriorityLow
	TaskPriorityMedium
	TaskPriorityHigh
)

type TaskState int

const (
	TaskStateIncomplete TaskState = iota // Default state
	TaskStateStarted
	TaskStateCompleted
)

type Task struct {
	ID           int64          `json:"id" db:"id"`
	Title        string         `json:"title" db:"title"`
	Description  *string        `json:"description" db:"description"`
	Due          *string        `json:"due" db:"dueUtc"`
	Priority     TaskPriority   `json:"priority" db:"priority"`
	CreatedUtc   string         `json:"createdUtc" db:"createdAtUtc"`
	UpdatedAtUtc *string        `json:"updatedAtUtc" db:"updatedAtUtc"`
	CompletedUtc *string        `json:"completedUtc" db:"completedAtUtc"`
	StartedUtc   sql.NullString `json:"startedUtc" db:"startedAtUtc"`
	State        TaskState      `json:"state" db:"state"`
}

func (task *Task) IsStarted() bool {
	log.Info().Str("startedAt", task.StartedUtc.String).Msg("IsStarted")
	return task.StartedUtc.Valid
}

func (task *Task) TimeSinceStarted() time.Duration {
	if !task.IsStarted() {
		return 0
	}
	var startedAt, err = time.Parse("2006-01-02 15:04:05", task.StartedUtc.String)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse startedAt")
		return 0
	}
	return time.Since(startedAt)
}

func (task *Task) PrettyAge(duration time.Duration) string {
	log.Info().Str("duration", duration.String()).Str("title", task.Title).Msg("PrettyAge")
	if duration.Hours() == 0 {
		duration = duration.Round(time.Minute)
	}
	if duration.Hours() > 24 {
		duration = duration.Truncate(time.Hour)
	}
	if duration.Minutes() < 1 {
		return "0m"
	}
	var pretty = duration.Truncate(time.Minute).String()
	return strings.TrimSuffix(pretty, "0s")
}

func (task *Task) TimeSinceStartedStr() string {
	return task.PrettyAge(task.TimeSinceStarted())
}

func (task *Task) TaskAge() time.Duration {
	var createdAt, err = time.Parse("2006-01-02 15:04:05", task.CreatedUtc)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse createdAt")
		return 0
	}
	return time.Since(createdAt)
}

func (task *Task) AgeStr() string {
	return task.PrettyAge(task.TaskAge())
}

func (store *Store) ListTasks() ([]Task, error) {
	var sql = ` SELECT * FROM tasks WHERE state != ? `
	var tasks []Task
	err := store.Con.Select(&tasks, sql, TaskStateCompleted)
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

func (store *Store) StartTimeToggleById(id int64) (*Task, error) {
	var sql = `
	UPDATE tasks
	SET
		updatedAtUtc = current_timestamp,
		startedAtUtc = case when startedAtUtc is null then current_timestamp else null end
	WHERE id = ?
	RETURNING *
	`
	var task = &Task{}
	var row = store.Con.QueryRowx(sql, id)
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
			priority, createdAtUtc, state,
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
		task.Priority, time.Now().UTC().String(), task.State,
		time.Now().UTC().String(), task.CompletedUtc, task.StartedUtc,
	)
	var err = row.StructScan(newTask)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (store *Store) CompletTaskById(taskId int64) (bool, error) {
	var sql = `
	UPDATE tasks
	SET
		state = ?,
		completedAtUtc = case
			when completedAtUtc is null
			then current_timestamp
			else completedAtUtc
		end
	WHERE id = ?
	`
	var res, err = store.Con.Exec(sql, TaskStateCompleted, taskId)
	if err != nil {
		return false, err
	}

	var affected int64
	affected, err = res.RowsAffected()

	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
