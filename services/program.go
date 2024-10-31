package services

import "github.com/luke-goddard/taskninja/events"

func (handler *ServiceHandler) RunProgram(e *events.RunProgram) ([]*events.Event, error) {
	var _, _, err = handler.interpreter.Execute(e.Program)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
