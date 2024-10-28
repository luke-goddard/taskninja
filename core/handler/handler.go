package handler

import (
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/interpreter"
)

type EventHandler struct {
	db          *db.Store
	sendChannel chan *events.Event
	interpreter *interpreter.Interpreter
}

func NewEventHandler(db *db.Store) *EventHandler {
	return &EventHandler{
		db:          db,
		sendChannel: make(chan *events.Event),
	}
}

func (handler *EventHandler) Send(event *events.Event) {
	handler.sendChannel <- event
}

func (handler *EventHandler) Run() {
	go func() {
		for {
			select {
			case event := <-handler.sendChannel:
				go handler.handleEvent(event)
			}
		}
	}()
}

func (handler *EventHandler) handleEvent(event *events.Event) {
	switch event.Type {
	case events.EventAddTask:
		handler.handleAddTask(event)
	}
}
