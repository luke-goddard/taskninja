package events

type EventType string

const (
	EventError EventType = "Error"
	EventRunProgram EventType = "RunProgram"
)

type Event struct {
	Type EventType
	Data interface{}
}
