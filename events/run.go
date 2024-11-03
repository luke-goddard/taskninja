package events

type RunProgram struct {
	Program string
}

func NewRunProgramEvent(program string) *Event {
	return &Event{
		Type: EventRunProgram,
		Data: &RunProgram{
			Program: program,
		},
	}
}

func DecodeRunProgramEvent(e *Event) *RunProgram {
	return e.Data.(*RunProgram)
}
