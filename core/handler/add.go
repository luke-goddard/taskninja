package handler

import "github.com/luke-goddard/taskninja/events"

func (handler *EventHandler) handleAddTask(event *events.Event) {
	var program = "add " + event.Data.(string)
	var _, _, err = handler.interpreter.Execute(program)
	if err != nil {
		handler.Send(events.NewErrorEvent(err))
		return
	}
}
