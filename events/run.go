package events

// ============================================================================
// RUN PROGRAM
// ============================================================================

// EventRunProgram is the event type for running a program e.g "add 'task'"
type RunProgram struct{ Program string }

// DecodeRunProgramEvent will decode the event to run a program
func DecodeRunProgramEvent(e *Event) *RunProgram { return e.Data.(*RunProgram) }

// NewRunProgramEvent will create a new event to run a program
func NewRunProgramEvent(program string) *Event {
	return &Event{
		Type: EventRunProgram,
		Data: &RunProgram{
			Program: program,
		},
	}
}
