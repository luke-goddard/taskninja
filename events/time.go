package events

type StartTaskById struct{ ID int64 }
type StopTaskById struct{ ID int64 }

func DecodeStartTaskByIdEvent(e *Event) *StartTaskById { return e.Data.(*StartTaskById) }
func DecodeStopTaskByIdEvent(e *Event) *StopTaskById { return e.Data.(*StopTaskById) }

func NewStartTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventStartTaskById,
		Data: &StartTaskById{
			ID: id,
		},
	}
}

func NewStopTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventStopTaskById,
		Data: &StopTaskById{
			ID: id,
		},
	}
}
