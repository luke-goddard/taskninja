package services

import (
	"context"

	"github.com/luke-goddard/taskninja/db"
)

func (handler *ServiceHandler) CreateTask(task *db.Task) (*db.Task, error){
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.CreateTask(ctx, task)
}
