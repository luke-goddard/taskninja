package events


func NewAddTaskEvent(program string) *Event {
	return &Event{
		Type: EventAddTask,
		Data: program,
	}
}
