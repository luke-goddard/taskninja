package handler

import "github.com/luke-goddard/taskninja/events"

func (handler *EventHandler) runProgram(e *events.RunProgram) []*events.Event {
	var newEvents, err = handler.services.RunProgram(e)
	if err != nil {
		var errorEvent = events.NewErrorEvent(err)
		return []*events.Event{errorEvent}
	}
	return newEvents
}
