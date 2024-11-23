package events

type EventType string

const (
	EventError            EventType = "Error"
	EventRunProgram       EventType = "RunProgram"
	EventListTasks        EventType = "ListTasks"
	EventCompleteTaskById EventType = "CompleteTaskByID"
	EventDeleteTaskById EventType = "DeleteTaskByID"
	EventStartTaskById    EventType = "StartTask"
	EventListTaskResponse EventType = "ListTaskResponse"
	EventIncreasePriority EventType = "IncreaseTaskPriority"
	EventDecreasePriority EventType = "DecreaseTaskPriority"
)

type Event struct {
	Type EventType
	Data interface{}
}
