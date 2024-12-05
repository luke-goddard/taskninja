package events

import "github.com/luke-goddard/taskninja/db"

// ============================================================================
// LIST
// ============================================================================

// ListTasks is an event to list all tasks
type ListTasks struct{}

// DecodeListTasksEvent will decode the event to list all tasks
func DecodeListTasksEvent(e *Event) *ListTasks { return e.Data.(*ListTasks) }

// NewListTasksEvent will create a new event to list all tasks
func NewListTasksEvent() *Event {
	return &Event{
		Type: EventListTasks,
		Data: &ListTasks{},
	}
}

// ============================================================================
// LIST RESPONSE
// ============================================================================

// ListTasksResponse is the response to list all tasks event
type ListTasksResponse struct{ Tasks []db.TaskDetailed }

// DecodeListTasksResponseEvent will decode the event to list all tasks response
func DecodeListTasksResponseEvent(e *Event) *ListTasksResponse { return e.Data.(*ListTasksResponse) }

// NewListTasksResponse will create a new event to list all tasks response
func NewListTasksResponse(tasks []db.TaskDetailed) *Event {
	return &Event{
		Type: EventListTaskResponse,
		Data: &ListTasksResponse{Tasks: tasks},
	}
}
