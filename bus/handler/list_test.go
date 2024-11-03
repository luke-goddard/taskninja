package handler

import (
	"testing"

	"github.com/luke-goddard/taskninja/events"
	"github.com/stretchr/testify/assert"
)

func TestListHandler(t *testing.T) {
	var handler = newTestHandler()

	var res, err = handler.services.Store.Con.Exec("INSERT INTO tasks (id, title, description, completed) VALUES (1, 'title', 'description', 0)")
	assert.Nil(t, err)

	_, err = res.LastInsertId()
	assert.Nil(t, err)

	var e = events.NewListTasksEvent()
	handler.bus.Publish(e)
}
