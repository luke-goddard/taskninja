package events

type EventType string

const (
	EventAddTask  EventType = "AddTask"
	EventError EventType = "Error"
)

type Event struct {
	Type EventType
	Data interface{}
}
