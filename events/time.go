package events

// ============================================================================
// START TIME
// ============================================================================
type StartTaskById struct{ ID int64 }

func DecodeStartTaskByIdEvent(e *Event) *StartTaskById { return e.Data.(*StartTaskById) }
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
type StopTaskById struct{ ID int64 }

func DecodeStopTaskByIdEvent(e *Event) *StopTaskById { return e.Data.(*StopTaskById) }
func NewStopTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventStopTaskById,
		Data: &StopTaskById{
			ID: id,
		},
	}
}
