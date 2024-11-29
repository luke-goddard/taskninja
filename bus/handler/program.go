package handler

import (
	"github.com/luke-goddard/taskninja/events"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) runProgram(e *events.RunProgram) []*events.Event {
	var _, err = handler.services.RunProgram(e.Program)
	if err != nil {
		log.Error().Err(err).Msg("error running program")
		var errorEvent = events.NewErrorEvent(err)
		return []*events.Event{errorEvent}
	}
	return []*events.Event{events.NewListTasksEvent()}
}
