package db

import (
	"fmt"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/rs/zerolog/log"
)

// Migrations contains all the migrations that need to be run.
// Each migration is a SQL statement. The migrations are run in order.
var Migrations = []string{
	M000_TaskSchema,
	M001_TagSchema,
	M002_TaskTagsSchema,
	M003_TaskSchema,
	M004_TaskSchema,
	M005_TaskSchema,
	M006_ProjectSchema,
	M007_TaskProjectsSchema,
	M008_TimeTrackingSchema,
	M009_TaskSchema,
	M010_TaskDependenciesSchema,
	M011_TaskSchema,
	"PRAGMA foreign_keys = ON",
}

// RunMigrations runs all migrations that have not been run.
// The PRAGMA user_version is used to determine the current schema version.
// If the schema version is less than the migration index, the migration is run.
// If the schema version is greater than the migration index,
// the migration is skipped.
func (store *Store) RunMigrations() error {
	var version = store.SchemaVersion()
	for i, migration := range Migrations {
		if i <= version && version != 0 {
			continue
		}
		_, err := store.Con.Exec(migration)
		if err != nil {
			log.Error().Err(err).Msg(migration)
			return fmt.Errorf("failed to run migration (%d): %w", i, err)
		}
	}
	return nil
}

// SchemaVersion returns the current schema version.
// The schema version is stored in the PRAGMA user_version.
func (store *Store) SchemaVersion() int {
	var row, err = store.Con.Query("PRAGMA user_version")
	assert.Nil(err, "failed to get schema version")
	var version int
	row.Next()
	err = row.Scan(&version)
	row.Close()
	assert.Nil(err, "failed to scan schema version")
	log.Debug().Int("version", version).Msg("schema version")
	return version
}
