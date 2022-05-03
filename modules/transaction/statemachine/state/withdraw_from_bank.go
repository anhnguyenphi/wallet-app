package state

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	BankWithdraw interface {
		CreateWithdrawTransaction(ctx context.Context, cardID int64, amount string) (string, error)
	}

	withdrawFromBankStateHandler struct {
		bank BankWithdraw
	}
)

func (w *withdrawFromBankStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return &dbclient.NoOpTx{}, nil
}

func (w *withdrawFromBankStateHandler) Handle(ctx context.Context, tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	bankingTxID, err := w.bank.CreateWithdrawTransaction(ctx, trans.ToWalletID.Int64, trans.Amount)
	if err != nil {
		return types.FailStatus, nil
	}
	trans.ExTransactionID = sql.NullString{String: bankingTxID}
	return types.SuccessStatus, nil
}

