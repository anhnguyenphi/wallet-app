package state

import (
	"context"
	"errors"
	"ewallet/modules/bank/service"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

var ErrRetryVerifyingBankTx = errors.New("retry banking transaction")

type (
	BankTransactionVerifier interface {
		GetStatus(ctx context.Context, cardID int64, bankingTransactionID string) (string, error)
	}

	verifyBankingTrxStateHandler struct {
		bank BankTransactionVerifier
	}
)

func (w *verifyBankingTrxStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return &dbclient.NoOpTx{}, nil
}

func (w *verifyBankingTrxStateHandler) Handle(ctx context.Context, tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	cardID := trans.FromCardID
	if trans.Type == types.DepositType {
		cardID = trans.ToCardID
	}
	status, err := w.bank.GetStatus(ctx, cardID.Int64, trans.ExTransactionID.String)
	if err != nil || status == service.StatusFailed {
		return types.FailStatus, nil
	}
	if status == service.StatusSuccess {
		return types.SuccessStatus, nil
	}
	return "", ErrRetryVerifyingBankTx
}

