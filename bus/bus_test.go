package bus

import (
	"testing"

	"github.com/luke-goddard/taskninja/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type SubscriberMock struct {
	called bool
}

func (s *SubscriberMock) Notify(e *events.Event) {
	s.called = true
}

func TestBus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bus Suite")
}

var _ = Describe("Bus", func() {
	var bus *Bus

	BeforeEach(func() {
		bus = NewBus()
	})

	It("should be created", func() {
		Expect(bus).ToNot(BeNil())
	})

	It("should not have subscribers", func() {
		Expect(bus.HasSubscribers()).To(BeFalse())
	})

	It("should have subscribers", func() {
		bus.Subscribe(&SubscriberMock{})
		Expect(bus.HasSubscribers()).To(BeTrue())
	})

	It("should publish", func() {
		var subscriber = &SubscriberMock{}
		bus.Subscribe(subscriber)
		bus.Publish(events.NewListTasksEvent())
		Expect(subscriber.called).To(BeTrue())
	})
})
