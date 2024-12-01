package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const M001_TagSchema = `
CREATE TABLE IF NOT EXISTS tags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);
PRAGMA user_version = 1;
`

// Tag is a struct that represents a tag e.g "+work"
type Tag struct {
	ID   int    `json:"id" db:"id"`     // Unique identifier of the tag
	Name string `json:"name" db:"name"` // Name of the tag
}

// CreateTagTx will create a new in the database (this should not exist)
// the transaction is NOT rolled back on err
func (store *Store) CreateTagTx(name string, tx *sqlx.Tx) (int64, error) {
	var res, err = tx.Exec("INSERT INTO tags (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("Failed creating new tag: %w", err)
	}
	var id int64
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed getting last insert id when creating a new tag: %w", err)
	}
	return id, nil
}

// CreateTag will create a new tag in the database (this should not exist)
func (store *Store) CreateTagCtx(ctx context.Context, name string) (int64, error) {
	var res, err = store.Con.ExecContext(ctx, "INSERT INTO tags (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("Failed to create a new tag: %w", err)
	}
	var id int64
	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Failed to get last insert id when creating a new tag: %w", err)
	}
	return id, nil
}

// GetTagByNameTx will get a single row for the tag with then name specified
// NOTE: the transaction is not rolled back on error
func (store *Store) GetTagByNameTx(name string, tx *sqlx.Tx) (*Tag, error) {
	var tag Tag
	var err = tx.Get(&tag, "SELECT * FROM tags WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tag by name: %w", err)
	}
	return &tag, nil
}

// GetTagByName will get a single row for the tag with the name specified
func (store *Store) GetTagByName(name string) (*Tag, error) {
	var tag Tag
	var err = store.Con.Get("SELECT * FROM tags WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("Failed to get a tag by name: %w", err)
	}
	return &tag, nil
}
