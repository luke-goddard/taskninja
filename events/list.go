package events

import "github.com/luke-goddard/taskninja/db"

type ListTasks struct{}

func NewListTasksEvent() *Event {
	return &Event{
		Type: EventListTasks,
		Data: &ListTasks{},
	}
}

func DecodeListTasksEvent(e *Event) *ListTasks {
	return e.Data.(*ListTasks)
}

type ListTasksResponse struct {
	Tasks []db.Task
}

func NewListTasksResponse(tasks []db.Task) *Event {
	return &Event{
		Type: EventListTaskResponse,
		Data: &ListTasksResponse{Tasks: tasks},
	}
}

func DecodeListTasksResponseEvent(e *Event) *ListTasksResponse {
	return e.Data.(*ListTasksResponse)
}
