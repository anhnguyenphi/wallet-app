package state

import (
	"context"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	rejectedStateHandler struct {
	}
)

func (w *rejectedStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return &dbclient.NoOpTx{}, nil
}

func (w *rejectedStateHandler) Handle(ctx context.Context, tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	return types.CompleteStatus, nil
}
