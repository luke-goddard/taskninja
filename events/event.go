package events

type EventType string

const (
	EventError            EventType = "Error"
	EventRunProgram       EventType = "RunProgram"
	EventListTasks        EventType = "ListTasks"
	EventDeleteTaskById   EventType = "DeleteTask"
	EventStartTaskById   EventType = "StartTask"
	EventListTaskResponse EventType = "ListTaskResponse"
)

type Event struct {
	Type EventType
	Data interface{}
}
