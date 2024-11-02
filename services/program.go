package services

import "github.com/luke-goddard/taskninja/events"

func (handler *ServiceHandler) RunProgram(e *events.RunProgram) error {
	var sql, args, err = handler.interpreter.Execute(e.Program)
	if err != nil {
		return err
	}

	_, err =handler.store.Con.Exec(string(sql), args...)

	return err
}
