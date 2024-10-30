package handler

import (
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/services"
)

type EventHandler struct {
	bus *bus.Bus
	services *services.ServiceHandler
}

func NewEventHandler(services *services.ServiceHandler, bus *bus.Bus) *EventHandler {
	return &EventHandler{services: services, bus: bus}
}

func (handler *EventHandler) Notify(e *events.Event) {
	var newEvents = handler.handle(e)
	for _, newEvent := range newEvents {
		handler.bus.Publish(newEvent)
	}
}

func (handler *EventHandler) handle(e *events.Event) []*events.Event {
	switch e.Type {
	case events.EventRunProgram:
		return handler.runProgram(events.DecodeRunProgramEvent(e))
	}
	return nil
}
