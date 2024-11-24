package db

import "database/sql"

const M008_TimeTrackingSchema = `
CREATE TABLE IF NOT EXISTS taskTime (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	taskId INTEGER NOT NULL,
	startTimeUtc TEXT NOT NULL DEFAULT current_timestamp,
	endTimeUtc TEXT,
	totalTime TEXT,
	FOREIGN KEY(taskId) REFERENCES task(id)
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

func (store *Store) StartTask(id int64) (*TaskTime, error) {
	var sql = ` INSERT INTO taskTime (taskId) VALUES (?) RETURNING *;`
	var row = store.Con.QueryRowx(sql, id)
	var taskTime = &TaskTime{}
	var err = row.StructScan(taskTime)
	if err != nil {
		return nil, err
	}
	return taskTime, nil
}

func (store *Store) StopTask(id int64) (*TaskTime, error) {
	var sql = `
	UPDATE taskTime
	SET
		endTimeUtc = current_timestamp,
		totalTime = (julianDay(current_timestamp) - julianDay(startTimeUtc)) * 24 * 60 * 60
	WHERE
		taskId = ? AND endTimeUtc IS NULL RETURNING *;
	`
	var row = store.Con.QueryRowx(sql, id)
	var taskTime = &TaskTime{}
	var err = row.StructScan(taskTime)
	if err != nil {
		return nil, err
	}
	return taskTime, nil
}

func (store *Store) GetCumTime(taskId int64) (int64, error) {
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
	var row = store.Con.QueryRowx(sql, taskId)
	var totalTime int64
	var err = row.Scan(&totalTime)
	if err != nil {
		return 0, err
	}
	return totalTime, nil
}
