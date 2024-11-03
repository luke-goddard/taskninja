package handler

import (
	"github.com/luke-goddard/taskninja/events"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) listTasks(e *events.ListTasks) []*events.Event {
	var tasks, err = handler.services.ListTasks(e)
	if err != nil {
		log.Error().Err(err).Msg("error listing tasks")
		var errorEvent = events.NewErrorEvent(err)
		return []*events.Event{errorEvent}
	}
	var resp = events.NewListTasksResponse(tasks)
	return []*events.Event{resp}
}