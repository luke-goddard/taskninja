package events


// ============================================================================
// INCREASE PRIORITY
// ============================================================================
type IncreasePriority struct { ID int64 }
func DecodeIncreasePriorityEvent(e *Event) *IncreasePriority { return e.Data.(*IncreasePriority) }
func NewIncreasePriorityEvent(id int64) *Event {
	return &Event{
		Type: EventIncreasePriority,
		Data: &IncreasePriority{
			ID: id,
		},
	}
}

// ============================================================================
// DECREASE PRIORITY
// ============================================================================

type DecreasePriority struct { ID int64 }
func DecodeDecreasePriorityEvent(e *Event) *DecreasePriority { return e.Data.(*DecreasePriority) }
func NewDecreasePriorityEvent(id int64) *Event {
	return &Event{
		Type: EventDecreasePriority,
		Data: &DecreasePriority{
			ID: id,
		},
	}
}

