package handler

import (
	"github.com/luke-goddard/taskninja/events"
	"github.com/rs/zerolog/log"
)

func (handler *EventHandler) startTaskById(e *events.StartTask) []*events.Event {
	log.Debug().Interface("event", e).Msg("starting task")
	var events = make([]*events.Event, 0)
	return events
}
