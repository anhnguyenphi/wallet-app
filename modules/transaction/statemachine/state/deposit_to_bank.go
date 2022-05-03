package state

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	BankDeposit interface {
		CreateDepositTransaction(ctx context.Context, cardID int64, amount string) (string, error)
	}

	depositToBankStateHandler struct {
		bank BankDeposit
	}
)

func (w *depositToBankStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return &dbclient.NoOpTx{}, nil
}

func (w *depositToBankStateHandler) Handle(ctx context.Context, tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	bankingTxID, err := w.bank.CreateDepositTransaction(ctx, trans.FromCardID.Int64, trans.Amount)
	if err != nil {
		return types.FailStatus, nil
	}
	trans.ExTransactionID = sql.NullString{String: bankingTxID}
	return types.SuccessStatus, nil
}


