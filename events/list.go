package events

import "github.com/luke-goddard/taskninja/db"

// ============================================================================
// LIST
// ============================================================================
type ListTasks struct{}

func DecodeListTasksEvent(e *Event) *ListTasks { return e.Data.(*ListTasks) }
func NewListTasksEvent() *Event {
	return &Event{
		Type: EventListTasks,
		Data: &ListTasks{},
	}
}

// ============================================================================
// LIST RESPONSE
// ============================================================================
type ListTasksResponse struct{ Tasks []db.TaskDetailed }

func DecodeListTasksResponseEvent(e *Event) *ListTasksResponse { return e.Data.(*ListTasksResponse) }
func NewListTasksResponse(tasks []db.TaskDetailed) *Event {
	return &Event{
		Type: EventListTaskResponse,
		Data: &ListTasksResponse{Tasks: tasks},
	}
}
