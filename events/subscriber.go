package events

// Subscriber is an interface that must be implemented by any subscriber
type Subscriber interface {
	Notify(e *Event) // Notify is called when an event is published
}
