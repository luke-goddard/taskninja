package db

import (
	"fmt"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/rs/zerolog/log"
)

var Migrations = []string{
	M000_TaskSchema,
	M001_TagSchema,
	M002_TaskTagsSchema,
	M003_TaskSchema,
}

func (store *Store) RunMigrations() error {
	var version = store.SchemaVersion()
	for i, migration := range Migrations {
		if i <= version && version != 0{
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
