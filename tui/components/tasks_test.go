package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/bus/handler"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/services"
	"github.com/luke-goddard/taskninja/tui/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type SubscriberMock struct {
	events []events.Event
}

func (s *SubscriberMock) Notify(e *events.Event) {
	s.events = append(s.events, *e)
}

func (s *SubscriberMock) HasEventOfType(eventType events.EventType) bool {
	for _, e := range s.events {
		if e.Type == eventType {
			return true
		}
	}
	return false
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Table Suite")
}

func newTestHandler() *services.ServiceHandler {
	var store = db.NewInMemoryStore()
	var interpreter = interpreter.NewInterpreter()
	return services.NewServiceHandler(interpreter, store)
}

var _ = Describe("Task Table", func() {
	var table *TaskTable
	var sub *SubscriberMock
	var bus_ *bus.Bus
	var service *services.ServiceHandler
	var eventHandler *handler.EventHandler

	BeforeEach(func() {
		service = newTestHandler()
		bus_ = bus.NewBus()
		eventHandler = handler.NewEventHandler(service, bus_)
		bus_.Subscribe(eventHandler)
		sub = &SubscriberMock{}
		bus_.Subscribe(sub)
		table = NewTaskTable(
			lipgloss.NewStyle(),
			&utils.TerminalDimensions{
				Width:  100,
				Height: 100,
			},
			utils.NewTheme(),
			bus_,
		)
		bus_.Subscribe(table)
		table.Init()
	})

	It("should be created", func() {
		Expect(table).ToNot(BeNil())
	})

	Describe("When table has rows", func() {
		BeforeEach(func() {
			bus_.Publish(events.NewRunProgramEvent(`add "T2"`))
			bus_.Publish(events.NewRunProgramEvent(`add "T1"`)) // <- this is the first row
		})
		It("should have rows", func() {
			Expect(table.Table.Rows()).To(HaveLen(2))
		})
		It("should have a row with the correct title", func() {
			var row = table.GetCurrentRow()
			Expect(row.Title()).To(Equal("T1"))
		})
		It("Pressing j should move selected row down", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			var row = table.GetCurrentRow()
			Expect(row.Title()).To(Equal("T2"))
		})
		It("Pressing k should move selected row up", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // down
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}) // up
			var row = table.GetCurrentRow()
			Expect(row.Title()).To(Equal("T1"))
		})
		It("Pressing D should delete the selected row", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})
			Expect(sub.HasEventOfType(events.EventDeleteTaskById)).To(BeTrue())
		})
		It("Pressing d should complete the selected row", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(sub.HasEventOfType(events.EventCompleteTaskById)).To(BeTrue())
		})
		It("Pressing s should start the selected row", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(sub.HasEventOfType(events.EventStartTaskById)).To(BeTrue())
		})
		It("Pressing + should increase the priority of the selected row", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
			Expect(sub.HasEventOfType(events.EventIncreasePriority)).To(BeTrue())
		})
		It("Pressing - should decrease the priority of the selected row", func() {
			table.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
			Expect(sub.HasEventOfType(events.EventDecreasePriority)).To(BeTrue())
		})
		It("Should have a zero urgency", func() {
			var row = table.GetCurrentRow()
			Expect(row.Urgency()).To(Equal(0.0))
		})
	})

	Describe("When table has no rows", func() {
		It("should not have rows", func() {
			Expect(table.Table.Rows()).To(HaveLen(0))
		})
	})

	Describe("When it has a task that is started", func() {
		BeforeEach(func() {
			var t, err = service.Store.CreateTask(&db.Task{
				Title: "T1",
				State: db.TaskStateStarted,
			})
			Expect(err).To(BeNil())
			var tid = t.ID

			err = service.Store.StartTrackingTaskTime(tid)
			Expect(err).To(BeNil())
			bus_.Publish(events.NewListTasksEvent())
		})
		It("should show the task as started", func() {
			Expect(table.Table.Rows()).To(HaveLen(1))
			var row = table.GetCurrentRow()
			Expect(row.Started()).To(BeTrue())
			Expect(row[TableColumnID]).To(Equal("1-⏰"))
		})
		It("should have a none zero urgency", func() {
			Expect(table.Table.Rows()).To(HaveLen(1))
			var row = table.GetCurrentRow()
			Expect(row.Urgency()).To(Equal(20.0))
		})
	})

	Describe("When a task has a priority", func() {
		BeforeEach(func() {
			service.Store.CreateTask(&db.Task{Title: "T4", Priority: db.TaskPriorityHigh})
			service.Store.CreateTask(&db.Task{Title: "T3", Priority: db.TaskPriorityMedium})
			service.Store.CreateTask(&db.Task{Title: "T2", Priority: db.TaskPriorityLow})
			service.Store.CreateTask(&db.Task{Title: "T1", Priority: db.TaskPriorityNone})
			bus_.Publish(events.NewListTasksEvent())
		})
		It("should show the task as high priority", func() {
			Expect(table.Table.Rows()).To(HaveLen(4))
			var row = table.GetRowAtPos(0)
			Expect(row.PriorityStr()).To(Equal("High"))
		})
		It("should show the task as medium priority", func() {
			var row = table.GetRowAtPos(1)
			Expect(row.PriorityStr()).To(Equal("Medium"))
		})
		It("should show the task as low priority", func() {
			var row = table.GetRowAtPos(2)
			Expect(row.PriorityStr()).To(Equal("Low"))
		})
		It("should show the task as none priority", func() {
			var row = table.GetRowAtPos(3)
			Expect(row.PriorityStr()).To(Equal("None"))
		})
	})
})
