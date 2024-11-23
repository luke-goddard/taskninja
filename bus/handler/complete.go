package handler

import "github.com/luke-goddard/taskninja/events"

func (handler *EventHandler) completeTaskById(e *events.CompleteTaskById) []*events.Event {
	var affected, err = handler.services.CompleteTaskById(e.Id)
	if err != nil {
		return []*events.Event{events.NewErrorEvent(err)}
	}
	if affected {
		return []*events.Event{events.NewListTasksEvent()}
	}
	return nil
}

func (handler *EventHandler) deleteTaskById(e *events.DeleteTaskById) []*events.Event {
	var affected, err = handler.services.DeleteTaskById(e.Id)
	if err != nil {
		return []*events.Event{events.NewErrorEvent(err)}
	}
	if affected {
		return []*events.Event{events.NewListTasksEvent()}
	}
	return nil
}
