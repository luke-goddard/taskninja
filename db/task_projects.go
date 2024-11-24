package db

const M007_TaskProjectsSchema = `
CREATE TABLE IF NOT EXISTS taskProjects (
	taskId INTEGER NOT NULL,
	projectId INTEGER NOT NULL,
	FOREIGN KEY (taskId) REFERENCES tasks(id) ON DELETE CASCADE,
	FOREIGN KEY (projectId) REFERENCES projects(id) ON DELETE CASCADE,
	UNIQUE (taskId, projectId),
	PRIMARY KEY (taskId, projectId)
);
PRAGMA user_version = 6;
`
