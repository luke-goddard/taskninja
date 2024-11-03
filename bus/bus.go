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

func NewBus() *Bus {
	return &Bus{
		subscribers: make([]events.Subscriber, 0),
	}
}

func (b *Bus) Subscribe(s events.Subscriber) {
	assert.NotNil(s, "subscriber is nil")
	b.subscribers = append(b.subscribers, s)
}

func (b *Bus) Publish(e *events.Event) {
	assert.NotNil(e, "event is nil")
	assert.True(b.HasSubscribers(), "no subscribers")

	for _, s := range b.subscribers {
		s.Notify(e)
	}
}

func (b *Bus) HasSubscribers() bool {
	return len(b.subscribers) > 0
}
