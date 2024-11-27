package db

import (
	"context"
	"database/sql"
	"fmt"
)

const M008_TimeTrackingSchema = `
CREATE TABLE IF NOT EXISTS taskTime (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	taskId INTEGER NOT NULL,
	startTimeUtc TEXT NOT NULL DEFAULT current_timestamp,
	endTimeUtc TEXT,
	totalTime TEXT,
	FOREIGN KEY(taskId) REFERENCES tasks(id) ON DELETE CASCADE
);
PRAGMA user_version = 8;
`

type TaskTime struct {
	Id           int64          `db:"id"`
	TaskId       int64          `db:"taskId"`
	StartTimeUtc string         `db:"startTimeUtc"`
	EndTimeUtc   sql.NullString `db:"endTimeUtc"`
	TotalTime    sql.NullString `db:"totalTime"`
}

func (store *Store) StartTrackingTaskTime(ctx context.Context, taskId int64) error {
	// 1. Set the task state to started
	// 2. If there are no times for the task, insert a new time
	// 3. If there are times for the task, do not insert a new time

	var tx, err = store.Con.Beginx()
	if err != nil {
		return fmt.Errorf("error starting task: %w", err)
	}
	defer tx.Rollback()

	var sql = `UPDATE tasks SET state = ? WHERE id = ?;`
	_, err = tx.ExecContext(ctx, sql, TaskStateStarted, taskId)
	if err != nil {
		return fmt.Errorf("error updating task state while starting task: %w", err)
	}

	sql = `
	  INSERT INTO taskTime (taskId)
	  SELECT ?
	  WHERE NOT EXISTS (SELECT 1 FROM taskTime WHERE taskId = ? AND endTimeUtc IS NULL);
	`
	_, err = tx.ExecContext(ctx, sql, taskId, taskId)
	if err != nil {
		return fmt.Errorf("error inserting task time while starting task: %w", err)
	}

	return tx.Commit()
}

func (store *Store) StopTrackingTaskTime(ctx context.Context, id int64) error {
	var sql = `UPDATE tasks SET state = 0 WHERE id = ?;`
	var tx, err = store.Con.Beginx()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	_, err = tx.Exec(sql, id)
	if err != nil {
		return fmt.Errorf("error updating task state while stopping task: %w", err)
	}

	sql = `
	UPDATE taskTime
	SET
		endTimeUtc = current_timestamp,
		totalTime = (julianDay(current_timestamp) - julianDay(startTimeUtc)) * 24 * 60 * 60
	WHERE
		taskId = ? AND endTimeUtc IS NULL;
	`
	_, err = tx.ExecContext(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("error updating task time while stopping task: %w", err)
	}

	return tx.Commit()
}

func (store *Store) GetTaskTimes(ctx context.Context, taskId int64) ([]TaskTime, error) {
	var sql = `SELECT * FROM taskTime WHERE taskId = ?;`
	var rows, err = store.Con.QueryxContext(ctx, sql, taskId)
	if err != nil {
		return nil, err
	}
	var taskTimes = []TaskTime{}
	for rows.Next() {
		var taskTime = TaskTime{}
		err = rows.StructScan(&taskTime)
		if err != nil {
			return nil, err
		}
		taskTimes = append(taskTimes, taskTime)
	}
	return taskTimes, nil
}

func (store *Store) GetCumTime(ctx context.Context, taskId int64) (int64, error) {
	var sql = `
	SELECT
	    SUM(
		CASE
		    WHEN endTimeUtc IS NULL THEN
			(julianday(current_timestamp) - julianday(startTimeUtc)) * 24 * 60 * 60
		    ELSE totalTime
		END
	    ) AS cumulativeTime
	FROM taskTime
	WHERE taskId = ?;
	`
	var row = store.Con.QueryRowxContext(ctx, sql, taskId)
	var totalTime int64
	var err = row.Scan(&totalTime)
	if err != nil {
		return 0, err
	}
	return totalTime, nil
}
