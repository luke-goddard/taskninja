package events

// NewErrorEvent will create a new event that is used to report an error
func NewErrorEvent(error error) *Event {
	return &Event{
		Type: EventError,
		Data: error,
	}
}

// DecodeErrorEvent will decode the event to report an error
func DecodeErrorEvent(e *Event) error { return e.Data.(error) }
