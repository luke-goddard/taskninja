package events

// ============================================================================
// COMPLETE BY ID
// ============================================================================

// CompleteTaskById is an event that is used to complete a task by its ID
type CompleteTaskById struct{ Id int64 }

// CompleteTaskById is an event that is used to complete a task by its ID
func DecodeCompletedTaskById(e *Event) *CompleteTaskById { return e.Data.(*CompleteTaskById) }

// NewCompleteEvent will create a new event to complete a task by its ID
func NewCompleteEvent(id int64) *Event {
	return &Event{
		Type: EventCompleteTaskById,
		Data: &CompleteTaskById{Id: id},
	}
}

// ============================================================================
// DELETE BY ID
// ============================================================================

// DeleteTaskById is an event that is used to delete a task by its ID
type DeleteTaskById struct{ Id int64 }

// DecodeDeleteTaskByIdEvent will decode the event to delete a task by its ID
func DecodeDeleteTaskByIdEvent(e *Event) *DeleteTaskById { return e.Data.(*DeleteTaskById) }

// DeleteTaskById is an event that is used to delete a task by its ID
func NewDeleteTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventDeleteTaskById,
		Data: &DeleteTaskById{Id: id},
	}
}
