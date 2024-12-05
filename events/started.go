package events

// ============================================================================
// START
// ============================================================================

type StartTask struct{ Id int64 }

// DecodeStartTaskEvent will decode the event to start a new task
func DecodeStartTaskEvent(e *Event) *StartTask { return e.Data.(*StartTask) }

// Used to start a new task
func NewStartTaskEvent(id int64) *Event {
	return &Event{
		Type: EventStartTaskById,
		Data: &StartTask{
			Id: id,
		},
	}
}
