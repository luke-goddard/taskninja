package events

type DeleteTask struct {
	Id int // The ID of the task to delete
}

func NewDeleteTaskEvent(id int) *Event {
	return &Event{
		Type: EventDeleteTask,
		Data: &DeleteTask{Id: id},
	}
}

func DecodeDeleteTaskEvent(e *Event) *DeleteTask {
	return e.Data.(*DeleteTask)
}
