package services

import "context"

func (handler *ServiceHandler) CreateNewTag(id int64) error {
	var ctx, cancle = context.WithDeadline(context.Background(), handler.timeout())
	defer cancle()
	return handler.Store.CreateTagTx(ctx, id)
}
