package services

import "context"

// CreateTag will create a new tag in the database (this should not exist)
func (handler *ServiceHandler) CreateNewTag(name string) (int64, error) {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.CreateTagCtx(ctx, name)
}
