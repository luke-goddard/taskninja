package handler

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/services"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func newTestHandler() *EventHandler {
	var store = db.NewInMemoryStore()
	var bus = bus.NewBus()
	var interpreter = interpreter.NewInterpreter()
	var srv = services.NewServiceHandler(interpreter, store)
	var handler = NewEventHandler(srv, bus)
	bus.Subscribe(handler)
	return handler
}

func TestDeletehandler(t *testing.T) {
	var handler = newTestHandler()

	var task, err = handler.services.CreateTask(&db.Task{Title: "title"})
	assert.Nil(t, err)
	assert.NotNil(t, task)

	var deleted bool
	deleted, err = handler.services.DeleteTasks(task.ID)
	assert.Nil(t, err)
	assert.True(t, deleted)

	// COUNT
	var count int64
	var rows *sqlx.Rows
	rows, err = handler.services.Store.Con.Queryx("SELECT COUNT(*) FROM tasks WHERE id = ?", task.ID)
	assert.Nil(t, err)
	rows.Next()
	err = rows.Scan(&count)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestListHandler(t *testing.T) {
	var handler = newTestHandler()
	var incompleteTask, err = handler.services.CreateTask(&db.Task{Title: "title"})
	assert.Nil(t, err)
	assert.NotNil(t, incompleteTask)

	var completed = "2024-11-23 13:58:09"
	_, err = handler.services.CreateTask(&db.Task{
		Title:        "title",
		CompletedUtc: &completed,
		State:        db.TaskStateCompleted,
	})
	assert.Nil(t, err)
	assert.NotNil(t, incompleteTask)

	var tasks []db.Task
	tasks, err = handler.services.ListTasks()
	assert.Nil(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, incompleteTask.ID, tasks[0].ID)
	for _, task := range tasks {
		log.Info().Interface("task", task).Msg("task")
	}
}

func TestStartTaskHandler(t *testing.T) {
	var handler = newTestHandler()
	var task, err = handler.services.CreateTask(&db.Task{
		Title: "title",
	})
	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.Empty(t, task.StartedUtc)
	assert.Equal(t, "title", task.Title)

	task, err = handler.services.StartTimeToggleById(task.ID)
	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.NotNil(t, task.StartedUtc)

	t.Run("restarting-a-started-task", func(t *testing.T) {
		task, err = handler.services.StartTimeToggleById(task.ID)
		assert.Nil(t, err)
		assert.NotNil(t, task)
		assert.Empty(t, task.StartedUtc) // Should be empty
	})
}

func TestCompleteHandler(t *testing.T) {
	var handler = newTestHandler()
	var task, err = handler.services.CreateTask(&db.Task{
		Title: "title",
	})
	assert.Nil(t, err)
	assert.NotNil(t, task)

}
