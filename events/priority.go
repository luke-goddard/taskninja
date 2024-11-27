package events

import "github.com/luke-goddard/taskninja/db"

// ============================================================================
// INCREASE PRIORITY
// ============================================================================
type IncreasePriority struct{ ID int64 }

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

type DecreasePriority struct{ ID int64 }

func DecodeDecreasePriorityEvent(e *Event) *DecreasePriority { return e.Data.(*DecreasePriority) }
func NewDecreasePriorityEvent(id int64) *Event {
	return &Event{
		Type: EventDecreasePriority,
		Data: &DecreasePriority{
			ID: id,
		},
	}
}

// ============================================================================
// SET PRIORITY
// ============================================================================

type SetPriority struct {
	ID       int64
	Priority db.TaskPriority
}

func DecodeSetPriorityEvent(e *Event) *SetPriority { return e.Data.(*SetPriority) }
func NewSetPriorityEvent(id int64, priority db.TaskPriority) *Event {
	return &Event{
		Type: EventSetPriority,
		Data: &SetPriority{
			ID:       id,
			Priority: priority,
		},
	}
}
