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
	var handler = newTestHandler()

	var res, err = handler.services.Store.Con.Exec("INSERT INTO tasks (id, title, description, completed) VALUES (1, 'title', 'description', 0)")
	assert.Nil(t, err)

	var id int64
	id, err = res.LastInsertId()
	assert.Nil(t, err)

	var e = events.NewDeleteTaskEvent(id)
	handler.bus.Publish(e)

	// COUNT
	var count int64
	var rows *sqlx.Rows
	rows, err = handler.services.Store.Con.Queryx("SELECT COUNT(*) FROM tasks WHERE id = ?", id)
	assert.Nil(t, err)
	rows.Next()
	err = rows.Scan(&count)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}

func newTestHandler() *EventHandler {
	var store = db.NewInMemoryStore()
	var bus = bus.NewBus()
	var interpreter = interpreter.NewInterpreter()
	var srv = services.NewServiceHandler(interpreter, store)
	var handler = NewEventHandler(srv, bus)
	bus.Subscribe(handler)
	return handler
}