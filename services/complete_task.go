package services

import (
	"context"
	"fmt"
)

func (handler *ServiceHandler) CompleteTaskById(taskId int64) (bool, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	// TODO: Convert to transaction
	var err = handler.Store.StopTrackingTaskTime(ctx, taskId)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, fmt.Errorf("Error stopping task time: %v", err)
		}
	}
	// TODO: Convert to transaction
	err = handler.Store.DeleteDependenciesForCompletedTask(taskId)
	if err != nil {
		return false, fmt.Errorf("Error deleting dependencies for completed task: %v", err)
	}
	return handler.Store.CompleteTaskById(taskId)
}
