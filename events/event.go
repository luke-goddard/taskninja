package events

type EventType string // EventType is a type of event

const (
	EventError            EventType = "Error"                // Error event
	EventRunProgram       EventType = "RunProgram"           // RunProgram e.g 'add task'
	EventListTasks        EventType = "ListTasks"            // ListTasks event
	EventCompleteTaskById EventType = "CompleteTaskByID"     // Mark a task as complete
	EventTableFuzzySearch EventType = "TableFuzzySearch"     // Fuzzy search for a task
	EventDeleteTaskById   EventType = "DeleteTaskByID"       // Delete a task
	EventStartTaskById    EventType = "StartTask"            // Start a task
	EventStopTaskById     EventType = "StopTask"             // Stop a task
	EventListTaskResponse EventType = "ListTaskResponse"     // List tasks responses to be consumed by the UI
	EventIncreasePriority EventType = "IncreaseTaskPriority" // Increase the priority of a task
	EventDecreasePriority EventType = "DecreaseTaskPriority" // Decrease the priority of a task
	EventSetPriority      EventType = "SetTaskPriority"      // Set the priority of a task
)

type Event struct {
	Type EventType   // Type of event
	Data interface{} // Data associated with the event (payload)
}
