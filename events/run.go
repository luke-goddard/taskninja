package events

// ============================================================================
// RUN PROGRAM
// ============================================================================
type RunProgram struct{ Program string }

func DecodeRunProgramEvent(e *Event) *RunProgram { return e.Data.(*RunProgram) }

func NewRunProgramEvent(program string) *Event {
	return &Event{
		Type: EventRunProgram,
		Data: &RunProgram{
			Program: program,
		},
	}
}
