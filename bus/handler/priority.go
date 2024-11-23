package handler

import "github.com/luke-goddard/taskninja/events"

func (handler *EventHandler) increasePriority(e *events.IncreasePriority) []*events.Event {
	var newEvents = make([]*events.Event, 0)
	var applied, err = handler.services.IncreasePriority(e.ID)
	if err != nil {
		newEvents = append(newEvents, events.NewErrorEvent(err))
		return newEvents
	}
	if applied {
		newEvents = append(newEvents, events.NewListTasksEvent())
	}
	return newEvents
}

func (handler *EventHandler) decreasePriority(e *events.DecreasePriority) []*events.Event {
	var newEvents = make([]*events.Event, 0)
	var applied, err = handler.services.DecreasePriority(e.ID)
	if err != nil {
		newEvents = append(newEvents, events.NewErrorEvent(err))
		return newEvents
	}
	if applied {
		newEvents = append(newEvents, events.NewListTasksEvent())
	}
	return newEvents
}
