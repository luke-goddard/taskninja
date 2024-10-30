package events

type Subscriber interface {
	Notify(e *Event)
}
