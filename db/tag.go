package db

const M001_TagSchema = `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);
PRAGMA user_version = 1;
`

// Tag is a struct that represents a tag e.g "+work"
type Tag struct {
	ID   int    `json:"id" db:"id"` // Unique identifier of the tag
	Name string `json:"name" db:"name"` // Name of the tag
}
