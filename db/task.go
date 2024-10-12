package db

const TaskSchema = `
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT,
	completed INTEGER NOT NULL DEFAULT 0,
	priority INTEGER NOT NULL DEFAULT 0 CHECK (priority >= 0 AND priority <= 3),
	dueUtc TEXT,
	updatedAtUtc TEXT NOT NULL DEFAULT current_timestamp,
	createdAtUtc TEXT NOT NULL DEFAULT current_timestamp,
	completedAtUtc TEXT,
);`

type TaskPriority int

const (
	TaskPriorityNone TaskPriority = iota
	TaskPriorityLow
	TaskPriorityMedium
	TaskPriorityHigh
)

type Task struct {
	ID              int          `json:"id" db:"id"`
	Description     string       `json:"description" db:"description"`
	Due             string       `json:"due" db:"due"`
	Completed       bool         `json:"completed" db:"completed"`
	Priority        TaskPriority `json:"priority" db:"priority"`
	CreatedUtc      string       `json:"createdUtc" db:"createdUtc"`
	UpdatedAtUtc    string       `json:"updatedAtUtc" db:"updatedAtUtc"`
	CompletedUtc    string       `json:"completedUtc" db:"completedUtc"`
}
