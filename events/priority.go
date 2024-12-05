package events

import "github.com/luke-goddard/taskninja/db"

// ============================================================================
// INCREASE PRIORITY
// ============================================================================

// IncreasePriority is an event to increase the priority of a task
type IncreasePriority struct{ ID int64 }

// DecodeIncreasePriorityEvent will decode the event to increase the priority of a task
func DecodeIncreasePriorityEvent(e *Event) *IncreasePriority { return e.Data.(*IncreasePriority) }

// NewIncreasePriorityEvent will create a new event to increase the priority of a task
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

// DecreasePriority is an event to decrease the priority of a task
type DecreasePriority struct{ ID int64 }

// DecodeDecreasePriorityEvent will decode the event to decrease the priority of a task
func DecodeDecreasePriorityEvent(e *Event) *DecreasePriority { return e.Data.(*DecreasePriority) }

// NewDecreasePriorityEvent will create a new event to decrease the priority of a task
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
	ID       int64           // ID of the task
	Priority db.TaskPriority // Priority to set
}

// DecodeSetPriorityEvent will decode the event to set the priority of a task
func DecodeSetPriorityEvent(e *Event) *SetPriority { return e.Data.(*SetPriority) }

// NewSetPriorityEvent will create a new event to set the priority of a task
func NewSetPriorityEvent(id int64, priority db.TaskPriority) *Event {
	return &Event{
		Type: EventSetPriority,
		Data: &SetPriority{
			ID:       id,
			Priority: priority,
		},
	}
}
