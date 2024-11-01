package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type Store struct {
	Con *sqlx.DB
}

func NewStore(conf *config.SqlConnectionConfig) (*Store, error) {
	var dsn = conf.DSN()
	log.Debug().Str("dsn", dsn).Msg("connecting to database")
	var con, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	var store = &Store{Con: con}
	err = store.RunMigrations()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (store *Store) RunMigrations() error {
	for i, migration := range Migrations {
		_, err := store.Con.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to run migration (%d): %w", i, err)
		}
	}
	return nil
}

func (store *Store) Close() {
	assert.True(store.IsConnected(), "store is not connected")
	store.Con.Close()
}

func (store *Store) IsConnected() bool {
	return store.Con != nil
}
