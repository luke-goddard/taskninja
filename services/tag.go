package services

import (
	"context"

	"github.com/luke-goddard/taskninja/db"
)

// CreateTag will create a new tag in the database (this should not exist)
func (handler *ServiceHandler) TagCreate(name string) (int64, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.TagCreate(ctx, name)
}

// Used to list all of the tags
func (handler *ServiceHandler) TagList() ([]db.Tag, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.TagList(ctx)
}
