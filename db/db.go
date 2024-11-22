package db

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type Store struct {
	Con      *sqlx.DB
	writeMut sync.Mutex
}

func NewInMemoryStore() *Store {
	var con, err = sqlx.Connect("sqlite3", ":memory:")
	assert.Nil(err, "failed to connect to in-memory database")
	var store = &Store{Con: con}
	err = store.RunMigrations()
	assert.Nil(err, "failed to run migrations")
	return store
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
		log.Error().Err(err).Msg("failed to run migrations")
		return nil, err
	}
	return store, nil
}

func (store *Store) Close() {
	assert.True(store.IsConnected(), "store is not connected")
	var err = store.Con.Close()
	if err != nil {
		assert.Fail("failed to close database connection %w", err)
	}
}

func (store *Store) IsConnected() bool {
	return store.Con != nil
}

