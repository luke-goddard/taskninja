package events

// ============================================================================
// START TIME
// ============================================================================

// StartTaskById is an event that starts a task by its ID
type StartTaskById struct{ ID int64 }

// DecodeStartTaskByIdEvent will decode the event to start a task by its ID
func DecodeStartTaskByIdEvent(e *Event) *StartTaskById { return e.Data.(*StartTaskById) }

// NewStartTaskByIdEvent will create a new event to start a task by its ID
func NewStartTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventStartTaskById,
		Data: &StartTaskById{
			ID: id,
		},
	}
}

// ============================================================================
// STOP TIME
// ============================================================================

// StopTaskById is an event that stops a task by its ID
type StopTaskById struct{ ID int64 }

// DecodeStopTaskByIdEvent will decode the event to stop a task by its ID
func DecodeStopTaskByIdEvent(e *Event) *StopTaskById { return e.Data.(*StopTaskById) }

// NewStopTaskByIdEvent will create a new event to stop a task by its ID
func NewStopTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventStopTaskById,
		Data: &StopTaskById{
			ID: id,
		},
	}
}
