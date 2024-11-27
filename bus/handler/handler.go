package handler

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/services"
	"github.com/rs/zerolog/log"
)

type EventHandler struct {
	bus      *bus.Bus
	services *services.ServiceHandler
}

func NewEventHandler(services *services.ServiceHandler, bus *bus.Bus) *EventHandler {
	assert.NotNil(services, "services is nil")
	assert.NotNil(bus, "bus is nil")
	return &EventHandler{services: services, bus: bus}
}

func (handler *EventHandler) Notify(e *events.Event) {
	assert.NotNil(e, "event is nil")
	var newEvents = handler.handle(e)
	for _, newEvent := range newEvents {
		handler.bus.Publish(newEvent)
	}
}

func (handler *EventHandler) handle(e *events.Event) []*events.Event {
	log.Debug().Interface("event", e).Msg("handling event")
	switch e.Type {
	case events.EventRunProgram:
		return handler.runProgram(events.DecodeRunProgramEvent(e))
	case events.EventListTasks:
		return handler.listTasks()
	case events.EventCompleteTaskById:
		return handler.completeTaskById(events.DecodeCompletedTaskById(e))
	case events.EventDeleteTaskById:
		return handler.deleteTaskById(events.DecodeDeleteTaskByIdEvent(e))
	case events.EventStartTaskById:
		return handler.startTaskById(events.DecodeStartTaskEvent(e))
	case events.EventStopTaskById:
		return handler.stopTaskById(events.DecodeStopTaskByIdEvent(e))
	case events.EventIncreasePriority:
		return handler.increasePriority(events.DecodeIncreasePriorityEvent(e))
	case events.EventDecreasePriority:
		return handler.decreasePriority(events.DecodeDecreasePriorityEvent(e))
	case events.EventSetPriority:
		return handler.setPriority(events.DecodeSetPriorityEvent(e))
	}
	return nil
}
