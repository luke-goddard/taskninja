package handler

import (
	"testing"
	"time"

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
	deleted, err = handler.services.DeleteTaskById(task.ID)
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

	var tasks []db.TaskDetailed
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
	assert.Equal(t, "title", task.Title)

	err = handler.services.StartTimeToggleById(task.ID)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	err = handler.services.StopTimeToggleById(task.ID)
	assert.Nil(t, err)

	tasks, err := handler.services.Store.ListTasks()
	assert.Nil(t, err)
	assert.Len(t, tasks, 1)
	log.Info().Interface("tasks", tasks).Msg("tasks")

	var detailedTask = tasks[0]
	assert.True(t, detailedTask.CumulativeTime.Valid)
	log.Info().Interface("detailedTask", detailedTask).Msg("detailedTask")
	t.Fail()

	// assert.Nil(t, err)
	// assert.NotNil(t, task)
	// assert.NotNil(t, task.StartedUtc)
	//
	// t.Run("restarting-a-started-task", func(t *testing.T) {
	// 	task, err = handler.services.StartTimeToggleById(task.ID)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, task)
	// 	assert.Empty(t, task.StartedUtc) // Should be empty
	// })
}

func TestCompleteHandler(t *testing.T) {
	var handler = newTestHandler()
	var task, err = handler.services.CreateTask(&db.Task{
		Title: "title",
	})
	assert.Nil(t, err)
	assert.NotNil(t, task)

}

func TestIncDecPriority(t *testing.T) {
	var handler = newTestHandler()
	var task, err = handler.services.CreateTask(&db.Task{
		Title:    "title",
		Priority: db.TaskPriorityMedium,
	})
	assert.Nil(t, err)
	handler.services.IncreasePriority(task.ID) // High
	handler.services.IncreasePriority(task.ID) // High
	handler.services.DecreasePriority(task.ID) // Medium
	handler.services.DecreasePriority(task.ID) // Low
	handler.services.DecreasePriority(task.ID) // None
	handler.services.DecreasePriority(task.ID) // None
	handler.services.IncreasePriority(task.ID) // Low

	task = handler.services.Store.GetTaskByIdOrPanic(task.ID)
	assert.Equal(t, db.TaskPriorityLow, task.Priority)
}
