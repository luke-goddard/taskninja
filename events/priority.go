package events

type IncreasePriority struct { ID int64 }
type DecreasePriority struct { ID int64 }

func DecodeIncreasePriorityEvent(e *Event) *IncreasePriority { return e.Data.(*IncreasePriority) }
func DecodeDecreasePriorityEvent(e *Event) *DecreasePriority { return e.Data.(*DecreasePriority) }

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

