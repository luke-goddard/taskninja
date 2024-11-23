package events

type IncreasePriority struct {
	ID int64 // The ID of the task to increase the priority of
}

type DecreasePriority struct {
	ID int64 // The ID of the task to increase the priority of
}

func NewIncreasePriorityEvent(id int64) *Event {
	return &Event{
		Type: EventIncreasePriority,
		Data: &IncreasePriority{
			ID: id,
		},
	}
}

func NewDecreasePriorityEvent(id int64) *Event {
	return &Event{
		Type: EventDecreasePriority,
		Data: &DecreasePriority{
			ID: id,
		},
	}
}
