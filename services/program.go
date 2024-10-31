package services

import "github.com/luke-goddard/taskninja/events"

func (handler *ServiceHandler) RunProgram(e *events.RunProgram) ([]*events.Event, error) {
	var sql, args, err = handler.interpreter.Execute(e.Program)
	if err != nil {
		return nil, err
	}

	handler.store.Con.MustExec(string(sql), args...)

	return nil, nil
}
