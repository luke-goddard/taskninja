package events

type DeleteTaskById struct {
	Id int64 // The ID of the task to delete
}

func NewDeleteTaskEvent(id int64) *Event {
	return &Event{
		Type: EventDeleteTaskById,
		Data: &DeleteTaskById{Id: id},
	}
}

func DecodeDeleteTaskEvent(e *Event) *DeleteTaskById {
	return e.Data.(*DeleteTaskById)
}
