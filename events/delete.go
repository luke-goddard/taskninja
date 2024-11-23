package events

type CompleteTaskById struct { Id int64 }
func DecodeCompletedTaskById(e *Event) *CompleteTaskById { return e.Data.(*CompleteTaskById) }

func NewCompleteEvent(id int64) *Event {
	return &Event{
		Type: EventCompleteTaskById,
		Data: &CompleteTaskById{Id: id},
	}
}


type DeleteTaskById struct { Id int64 }
func DecodeDeleteTaskByIdEvent(e *Event) *DeleteTaskById { return e.Data.(*DeleteTaskById) }

func NewDeleteTaskByIdEvent(id int64) *Event {
	return &Event{
		Type: EventDeleteTaskById,
		Data: &DeleteTaskById{Id: id},
	}
}


