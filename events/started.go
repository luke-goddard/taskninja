package events

type StartTask struct {
	Id int64 // The ID of the task to start
}

// Used to start a new task
func NewStartTaskEvent(id int64) *Event {
	return &Event{
		Type: EventStartTaskById,
		Data: &StartTask{
			Id: id,
		},
	}
}

func DecodeStartTaskEvent(e *Event) *StartTask {
	return e.Data.(*StartTask)
}
