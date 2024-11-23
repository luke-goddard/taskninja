package handler

import (
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/events"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) startTaskById(e *events.StartTask) []*events.Event {
	log.Debug().Interface("event", e).Msg("starting task")
	var newEvents = make([]*events.Event, 0)
	var task, err = handler.services.StartTimeToggleById(e.Id)
	if err != nil {
		log.Error().Err(err).Msg("error starting task")
		newEvents = append(newEvents, events.NewErrorEvent(err))
		return newEvents
	}
	assert.NotNil(task, "task is nil, even though error is nil")
	newEvents = append(newEvents, events.NewListTasksEvent())
	return newEvents
}
