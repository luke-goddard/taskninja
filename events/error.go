package events

// ============================================================================
// DELETE BY ID
// ============================================================================
func NewErrorEvent(error error) *Event {
	return &Event{
		Type: EventError,
		Data: error,
	}
}
func DecodeErrorEvent(e *Event) error { return e.Data.(error) }
