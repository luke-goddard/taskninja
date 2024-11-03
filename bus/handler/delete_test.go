package handler

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/services"
	"github.com/stretchr/testify/assert"
)

func TestDeletehandler(t *testing.T) {
	var store = db.NewInMemoryStore()
	var bus = bus.NewBus()
	var interpreter = interpreter.NewInterpreter()
	var srv = services.NewServiceHandler(interpreter, store)
	var handler = NewEventHandler(srv, bus)
	bus.Subscribe(handler)

	var res, err = store.Con.Exec("INSERT INTO tasks (id, title, description, completed) VALUES (1, 'title', 'description', 0)")
	assert.Nil(t, err)

	var id int64
	id, err = res.LastInsertId()
	assert.Nil(t, err)

	var e = events.NewDeleteTaskEvent(id)
	bus.Publish(e)

	// COUNT
	var count int64
	var rows *sqlx.Rows
	rows, err = store.Con.Queryx("SELECT COUNT(*) FROM tasks WHERE id = ?", id)
	assert.Nil(t, err)
	rows.Next()
	err = rows.Scan(&count)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}
