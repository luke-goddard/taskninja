package handler

import (
	"github.com/luke-goddard/taskninja/events"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) startTaskById(e *events.StartTask) []*events.Event {
	log.Debug().Interface("event", e).Msg("starting task")
	var newEvents = make([]*events.Event, 0)
	var err = handler.services.StartTimeToggleById(e.Id)
	if err != nil {
		log.Error().Err(err).Msg("error starting task")
		newEvents = append(newEvents, events.NewErrorEvent(err))
		return newEvents
	}
	newEvents = append(newEvents, events.NewListTasksEvent())
	return newEvents
}

func (handler *EventHandler) stopTaskById(e *events.StopTaskById) []*events.Event {
	log.Debug().Interface("event", e).Msg("stopping task")
	var newEvents = make([]*events.Event, 0)
	var err = handler.services.StopTimeToggleById(e.ID)
	if err != nil {
		log.Error().Err(err).Msg("error stopping task")
		newEvents = append(newEvents, events.NewErrorEvent(err))
		return newEvents
	}
	newEvents = append(newEvents, events.NewListTasksEvent())
	return newEvents
}
