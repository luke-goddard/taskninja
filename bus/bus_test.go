package bus

import (
	"testing"

	"github.com/luke-goddard/taskninja/events"
	"github.com/stretchr/testify/assert"
)

type SubscriberMock struct {
	called bool
}

func (s *SubscriberMock) Notify(e *events.Event) {
	s.called = true
}

func TestBus(t *testing.T) {
	var bus = NewBus()
	assert.NotNil(t, bus, "bus is nil")
	assert.False(t, bus.HasSubscribers(), "bus has subscribers")

	var subscriber = &SubscriberMock{}
	assert.False(t, subscriber.called, "subscriber has been called")

	bus.Subscribe(subscriber)
	assert.True(t, bus.HasSubscribers(), "bus has no subscribers")

	bus.Publish(events.NewListTasksEvent())
	assert.True(t, subscriber.called, "subscriber has not been called")
}


