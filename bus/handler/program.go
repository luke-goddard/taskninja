package handler

import (
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) runProgram(e *events.RunProgram) []*events.Event {
	var program, err = handler.services.RunProgram(e.Program)
	if err != nil {
		log.Error().Err(err).Msg("error running program")
		var errorEvent = events.NewErrorEvent(err)
		return []*events.Event{errorEvent}
	}
	if program.Kind == ast.CommandKindAdd {
		return []*events.Event{events.NewListTasksEvent()}
	}
	return nil
}
