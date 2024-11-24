package events

import (
	"testing"

	"github.com/luke-goddard/taskninja/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Events Suite")
}

// ============================================================================
// DELETE.go
// ============================================================================
var _ = Describe("CompleteTaskById", func() {
	var event = NewCompleteEvent(1)

	It("should decode", func() {
		Expect(DecodeCompletedTaskById(event).Id).To(Equal(int64(1)))
	})
})

var _ = Describe("DeleteTaskById", func() {
	var event = NewDeleteTaskByIdEvent(1)

	It("should decode", func() {
		Expect(DecodeDeleteTaskByIdEvent(event).Id).To(Equal(int64(1)))
	})
})

// ============================================================================
// ERROR.go
// ============================================================================
var _ = Describe("NewErrorEvent", func() {
	var event = NewErrorEvent(nil)

	It("should create", func() {
		Expect(event.Type).To(Equal(EventError))
		Expect(event.Data).To(BeNil())
	})
})

// ============================================================================
// LIST.go
// ============================================================================
var _ = Describe("NewListEvent", func() {
	var event = NewListTasksEvent()

	It("should create", func() {
		Expect(event.Type).To(Equal(EventListTasks))
		Expect(event.Data).To(Equal(&ListTasks{}))
	})
})

var _ = Describe("NewListResponseEvent", func() {
	var event = NewListTasksResponse([]db.TaskDetailed{
		{Task: db.Task{ID: 1, Title: "do the dishes"}},
	})

	It("should create", func() {
		Expect(event.Type).To(Equal(EventListTaskResponse))
		Expect(event.Data).To(Equal(&ListTasksResponse{Tasks: []db.TaskDetailed{
			{Task: db.Task{ID: 1, Title: "do the dishes"}},
		}}))
	})
})

// ============================================================================
// PRIORITY.go
// ============================================================================

var _ = Describe("IncreasePriority", func() {
	var event = NewIncreasePriorityEvent(1)

	It("should decode", func() {
		Expect(DecodeIncreasePriorityEvent(event).ID).To(Equal(int64(1)))
	})
})

var _ = Describe("DecreasePriority", func() {
	var event = NewDecreasePriorityEvent(1)

	It("should decode", func() {
		Expect(DecodeDecreasePriorityEvent(event).ID).To(Equal(int64(1)))
	})
})

// ============================================================================
// RUN.go
// ============================================================================
var _ = Describe("Run", func() {
	var event = NewRunProgramEvent("hehe")

	It("should decode", func() {
		Expect(DecodeRunProgramEvent(event).Program).To(Equal("hehe"))
	})
})

// ============================================================================
// START.go
// ============================================================================
var _ = Describe("Start", func() {
	var event = NewStartTaskEvent(1)

	It("should decode", func() {
		Expect(DecodeStartTaskEvent(event).Id).To(Equal(int64(1)))
	})
})

// ============================================================================
// TIME.go
// ============================================================================
var _ = Describe("StartTask", func() {
	var event = NewStartTaskByIdEvent(1)

	It("should decode", func() {
		Expect(DecodeStartTaskByIdEvent(event).ID).To(Equal(int64(1)))
	})
})

var _ = Describe("StopTask", func() {
	var event = NewStopTaskByIdEvent(1)

	It("should decode", func() {
		Expect(DecodeStopTaskByIdEvent(event).ID).To(Equal(int64(1)))
	})
})
