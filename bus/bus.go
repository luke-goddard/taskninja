package bus

import "github.com/luke-goddard/taskninja/events"

// FANOUT pubsub event bus
// honk honk
type Bus struct {
	subscribers []events.Subscriber
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make([]events.Subscriber, 0),
	}
}

func (b *Bus) Subscribe(s events.Subscriber) {
	b.subscribers = append(b.subscribers, s)
}

func (b *Bus) Publish(e *events.Event) {
	for _, s := range b.subscribers {
		go s.Notify(e)
	}
}

