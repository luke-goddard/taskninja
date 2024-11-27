package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/luke-goddard/taskninja/assert"
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

const M009_TaskSchema = `
ALTER TABLE tasks DROP COLUMN startedAtUtc;
PRAGMA user_version = 9;
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

type UrgencyCoefficient float64

const URGENCY_MAX_AGES = time.Duration(365 * 24 * time.Hour)
const (
	URGENCY_NEXT_TAG_COEFFICIENT        UrgencyCoefficient = 15.0 // +Next
	URGENCY_PRIORITY_HIGH_COEFFICIENT   UrgencyCoefficient = 4.0  // P:High
	URGENCY_PRIORITY_MEDIUM_COEFFICIENT UrgencyCoefficient = 2.0  // P:Med
	URGENCY_PRIORITY_LOW_COEFFICIENT    UrgencyCoefficient = 1.0  // P:Low
	URGENCY_PRIORITY_NONE_COEFFICENT    UrgencyCoefficient = 0.0  // P:None
	URGENCY_DUE_COEFFICIENT             UrgencyCoefficient = 12.0 // Due:now
	URGENCY_BLOCKING_COEFFICIENT        UrgencyCoefficient = 8.0  // Task Dependencies
	URGENCY_ACTIVE_COEFFICIENT          UrgencyCoefficient = 20.0 // Task is started
	URGENCY_SCHEDULED_COEFFICIENT       UrgencyCoefficient = 5.0  // Task is scheduled
	URGENCY_AGE_COEFFICIENT             UrgencyCoefficient = 2.0  // Task age
	URGENCY_ANNOTATIONS_COEFFICIENT     UrgencyCoefficient = 1.0  // Task has annotations
	URGENCY_TAGS_COEFFICIENT            UrgencyCoefficient = 1.0  // Task has tags
	URGENCY_PROJECT_COEFFICIENT         UrgencyCoefficient = 1.0  // Task has project
	URGENCY_BLOCKED_COEFFICIENT         UrgencyCoefficient = -5.0 // Task is blocked
	URGENCY_WAITING_COEFFICIENT         UrgencyCoefficient = -3.0 // Task is waiting
)

const EPSILION = 0.000001
const SQLITE_TIME_FORMAT = "2006-01-02 15:04:05"

type Task struct {
	ID           int64          `json:"id" db:"id"`
	Title        string         `json:"title" db:"title"`
	Priority     TaskPriority   `json:"priority" db:"priority"`
	CreatedUtc   string         `json:"createdUtc" db:"createdAtUtc"`
	State        TaskState      `json:"state" db:"state"`
	Description  sql.NullString `json:"description" db:"description"`
	Due          sql.NullString `json:"due" db:"dueUtc"`
	UpdatedAtUtc sql.NullString `json:"updatedAtUtc" db:"updatedAtUtc"`
	CompletedUtc sql.NullString `json:"completedUtc" db:"completedAtUtc"`
}

type TaskDetailed struct {
	Task

	ProjectCount    int            `json:"projectCount" db:"projectCount"`
	ProjectNames    sql.NullString `json:"projectNames" db:"projectNames"`
	FirstStartedUtc sql.NullString `json:"firstStartedUtc" db:"firstStartedUtc"`
	CumulativeTime  sql.NullString `json:"cumulativeTime" db:"cumulativeTime"`
	Inprogress      bool           `json:"inprogress" db:"inprogress"`
	urgencyComputed float64
}

func (task *Task) PriorityStr() string {
	switch task.Priority {
	case TaskPriorityLow:
		return "Low"
	case TaskPriorityMedium:
		return "Medium"
	case TaskPriorityHigh:
		return "High"
	default:
		return "None"
	}
}

func (task *Task) PrettyAge(duration time.Duration) string {
	if duration.Hours() == 0 {
		duration = duration.Round(time.Minute)
	}

	if duration.Hours() >= 24*7 {
		var weeks = duration.Hours() / (24 * 7)
		var days = int(duration.Hours()) % (24 * 7) / 7
		return fmt.Sprintf("%dw%dd", int(weeks), int(days))
	}

	if duration.Hours() >= 24 {
		var days = duration.Hours() / 24
		var hours = int(duration.Hours()) % 24
		return fmt.Sprintf("%dd%dh", int(days), int(hours))
	}

	if duration.Minutes() < 1 {
		return "0m"
	}
	var pretty = duration.Truncate(time.Minute).String()
	return strings.TrimSuffix(pretty, "0s")
}

func (task *Task) AgeTime() time.Duration {
	var createdAt, err = time.Parse(SQLITE_TIME_FORMAT, task.CreatedUtc)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse createdAt")
		return 0
	}
	return time.Since(createdAt)
}

func (task *Task) AgeStr() string {
	return task.PrettyAge(task.AgeTime())
}

func (task *TaskDetailed) IsStarted() bool {
	return task.Inprogress
}

func (task *TaskDetailed) UrgencyStr() string {
	return fmt.Sprintf("%.2f", task.Urgency())
}

func (task *TaskDetailed) PrettyCumTime() string {
	if !task.CumulativeTime.Valid {
		return ""
	}
	var duration, err = time.ParseDuration(task.CumulativeTime.String + "s")
	if err != nil {
		log.Error().Err(err).Msg("failed to parse cumulative time")
		return ""
	}
	return task.PrettyAge(duration)
}

func (task *TaskDetailed) Urgency() float64 {
	if task.urgencyComputed == 0.0 {
		task.urgencyComputed = task.urgency()
	}
	return task.urgencyComputed
}

func (task *TaskDetailed) urgency() float64 {
	return task.urgencyProject() +
		task.urgencyActive() +
		task.urgencyScheduled() +
		task.urgencyDue() +
		task.urgencyAge() +
		task.urgencyPriority()

	// TODO add these when we have them
	// task.urgencyWaiting() +
	// task.urgencyBlocked() +
	// task.urgencyBlocking() +
	// task.urgencyTags() +
	// task.urgencyAnnotations()
}

func (task *TaskDetailed) urgencyProject() float64 {
	if URGENCY_PROJECT_COEFFICIENT < EPSILION || task.ProjectCount == 0 {
		return 0
	}
	return float64(URGENCY_PROJECT_COEFFICIENT)
}

func (task *TaskDetailed) urgencyActive() float64 {
	if URGENCY_ACTIVE_COEFFICIENT < EPSILION || !task.Inprogress {
		return 0
	}
	return float64(URGENCY_ACTIVE_COEFFICIENT)
}

func (task *TaskDetailed) urgencyScheduled() float64 {
	if URGENCY_SCHEDULED_COEFFICIENT < EPSILION || task.Due.Valid == false {
		return 0.0
	}
	return float64(URGENCY_SCHEDULED_COEFFICIENT)
}

func (task *TaskDetailed) urgencyAge() float64 {
	if URGENCY_DUE_COEFFICIENT < EPSILION {
		return 0.0
	}
	var age = task.AgeTime()
	var ageDays = age.Hours() / 24

	if ageDays > URGENCY_MAX_AGES.Hours()/24 {
		return 1.0
	}
	return 1.0 * ageDays / URGENCY_MAX_AGES.Hours() / 24
}

func (task *TaskDetailed) urgencyPriority() float64 {
	switch task.Priority {
	case TaskPriorityHigh:
		return float64(URGENCY_PRIORITY_HIGH_COEFFICIENT)
	case TaskPriorityMedium:
		return float64(URGENCY_PRIORITY_MEDIUM_COEFFICIENT)
	case TaskPriorityLow:
		return float64(URGENCY_PRIORITY_LOW_COEFFICIENT)
	default:
		return float64(URGENCY_PRIORITY_NONE_COEFFICENT)
	}
}

//	Past                  Present                              Future
//	Overdue               Due                                     Due
//
//	-7 -6 -5 -4 -3 -2 -1  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 days
//
// <-- 1.0                         linear                            0.2 -->
//
//	capped                                                        capped
//
// Ported from https://github.com/GothenburgBitFactory/taskwarrior/blob/develop/src/Task.cpp#L1702
func (task *TaskDetailed) urgencyDue() float64 {
	if URGENCY_DUE_COEFFICIENT < EPSILION || task.Due.Valid == false {
		return 0.0
	}
	var due, err = time.Parse(SQLITE_TIME_FORMAT, *&task.Due.String)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse due")
		return 0.0
	}
	var duration = time.Since(due)
	var daysOverDue = duration.Hours() / 24

	if daysOverDue > 7 {
		return 1.0
	} else if daysOverDue >= -14 {
		return ((daysOverDue + 14.0) * 0.8 / 21.0) + 0.2
	} else {
		return 0.2
	}
}

func (store *Store) CountTasks() (int64, error) {
	var sql = `SELECT COUNT(*) FROM tasks`
	var count int64
	err := store.Con.Get(&count, sql)
	if err != nil {
		log.Error().Err(err).Msg("failed to count tasks")
		return -1, err
	}
	return count, nil
}
func (store *Store) ListTasks() ([]TaskDetailed, error) {
	var sql = `
	SELECT
		tasks.*,
		COUNT(taskProjects.projectId) AS projectCount,
		GROUP_CONCAT(projects.title ORDER BY projects.title ASC) AS projectNames,
		MIN(taskTime.startTimeUtc) AS firstStartedUtc,
		CASE
			WHEN SUM(
				CASE
					WHEN taskTime.endTimeUtc IS NULL
					AND taskTime.startTimeUtc IS NOT NULL
					THEN 1
					ELSE 0
					END) > 0 THEN 1
			ELSE 0
		END AS inprogress,
		SUM(
			CASE
			    WHEN taskTime.endTimeUtc IS NULL THEN
				(julianday(current_timestamp) - julianday(taskTime.startTimeUtc)) * 24 * 60 * 60
			    ELSE taskTime.totalTime
			END
		) AS cumulativeTime
	FROM tasks
	LEFT JOIN taskProjects ON taskProjects.taskId = tasks.id
	LEFT JOIN projects ON projects.id = taskProjects.projectId
	LEFT JOIN taskTime ON taskTime.taskId = tasks.id
	WHERE
		tasks.state != 2 -- COMPLETED
	GROUP BY tasks.id;
	`
	var tasks []TaskDetailed
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

func (store *Store) CreateTask(task *Task) (*Task, error) {
	var sql = `
	INSERT INTO tasks
		(
			title, description, dueUtc,
			priority, createdAtUtc, state,
			updatedAtUtc, completedAtUtc
		)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING *
	`
	var newTask = &Task{}
	var row = store.Con.QueryRowx(
		sql,
		task.Title, task.Description, task.Due,
		task.Priority, time.Now().UTC().String(), task.State,
		time.Now().UTC().String(), task.CompletedUtc,
	)
	var err = row.StructScan(newTask)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (store *Store) CompleteTaskById(taskId int64) (bool, error) {
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

func (store *Store) IncreasePriority(id int64) (bool, error) {
	var sql = `
	UPDATE tasks
	SET
		priority = case
			when priority = 0 then 1
			when priority = 1 then 2
			when priority = 2 then 3
			when priority = 3 then 3
			else 0
		end
	WHERE id = ?
	`
	var res, err = store.Con.Exec(sql, id)
	if err != nil {
		return false, err
	}

	var affected int64
	affected, err = res.RowsAffected()
	assert.True(affected <= 1, "affected should be 0 or 1")
	return affected == 1, err
}

func (store *Store) DecreasePriority(id int64) (bool, error) {
	var sql = `
	UPDATE tasks
	SET
		priority = case
			when priority = 0 then 0
			when priority = 1 then 0
			when priority = 2 then 1
			when priority = 3 then 2
			else 0
		end
	WHERE id = ?
	`
	var res, err = store.Con.Exec(sql, id)
	if err != nil {
		return false, err
	}

	var affected int64
	affected, err = res.RowsAffected()
	assert.True(affected <= 1, "affected should be 0 or 1")
	return affected == 1, err
}

func (store *Store) GetTaskByIdOrPanic(id int64) *Task {
	var sql = `SELECT * FROM tasks WHERE id = ?`
	var task = &Task{}
	var err = store.Con.Get(task, sql, id)
	if err != nil {
		assert.Nil(err, "failed to get task by id")
	}
	return task
}

func (store *Store) SetTaskStateToStarted(taskId int64) error {
	return store.SetTaskState(taskId, TaskStateStarted)
}

func (store *Store) SetTaskStateToIncomplete(taskId int64) error {
	return store.SetTaskState(taskId, TaskStateIncomplete)
}

func (store *Store) SetTaskState(taskId int64, state TaskState) error {
	var sql = `UPDATE tasks SET state = ? WHERE id = ?; `
	_, err := store.Con.Exec(sql, state, taskId)
	return err
}

func (store *Store) GetTaskById(taskId int64) (*Task, error) {
	var sql = `SELECT * FROM tasks WHERE id = ?`
	var task = &Task{}
	var err = store.Con.Get(task, sql, taskId)
	if err != nil {
		return nil, err
	}
	return task, err
}
