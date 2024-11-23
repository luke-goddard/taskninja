package events

type CompleteTaskById struct {
	Id int64 // The ID of the task to delete
}

func NewCompleteEvent(id int64) *Event {
	return &Event{
		Type: EventCompleteTaskById,
		Data: &CompleteTaskById{Id: id},
	}
}

func DecodeDeleteTaskEvent(e *Event) *CompleteTaskById {
	return e.Data.(*CompleteTaskById)
}
