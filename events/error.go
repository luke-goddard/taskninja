package events

func NewErrorEvent(error error) *Event {
	return &Event{
		Type: EventError,
		Data: error,
	}
}
