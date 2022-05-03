package state

import (
	"context"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	completeStateHandler struct {
	}
)

func (w *completeStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return &dbclient.NoOpTx{}, nil
}

func (w *completeStateHandler) Handle(ctx context.Context, tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	return types.CompleteStatus, nil
}


