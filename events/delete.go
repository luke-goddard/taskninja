package events

type DeleteTaskById struct {
	Id int // The ID of the task to delete
}

func NewDeleteTaskEvent(id int) *Event {
	return &Event{
		Type: EventDeleteTaskById,
		Data: &DeleteTaskById{Id: id},
	}
}

func DecodeDeleteTaskEvent(e *Event) *DeleteTaskById {
	return e.Data.(*DeleteTaskById)
}
