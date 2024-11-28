package bus

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/events"

)

// FANOUT pubsub event bus
// honk honk
type Bus struct {
	subscribers []events.Subscriber
}

// NewBus creates a new event bus, don't forget to subscribe
func NewBus() *Bus {
	return &Bus{
		subscribers: make([]events.Subscriber, 0),
	}
}

// Subscribe to the event bus so that you can handle events
func (b *Bus) Subscribe(s events.Subscriber) {
	assert.NotNil(s, "subscriber is nil")
	b.subscribers = append(b.subscribers, s)
}

// Publish an event to all subscribers
func (b *Bus) Publish(e *events.Event) {
	assert.NotNil(e, "event is nil")
	assert.True(b.HasSubscribers(), "no subscribers")

	for _, s := range b.subscribers {
		s.Notify(e)
	}
}

// HasSubscribers returns true if there are subscribers
func (b *Bus) HasSubscribers() bool {
	return len(b.subscribers) > 0
}
