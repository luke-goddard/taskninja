package db

const M006_ProjectSchema = `
CREATE TABLE IF NOT EXISTS projects (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL
);
PRAGMA user_version = 6;
`
