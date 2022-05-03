package state

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	validatingStateHandler struct {
		card dbclient.DB
	}
)

func (w *validatingStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return w.card.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (v *validatingStateHandler) Handle(ctx context.Context,tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	if trans.Type == types.DepositType || trans.Type == types.WithdrawType {
		cardID, walletID := trans.ToCardID, trans.FromWalletID
		if trans.Type == types.DepositType {
			cardID, walletID = trans.FromCardID, trans.ToWalletID
		}
		ok, err := v.verifyCard(ctx, tx, cardID.Int64, walletID.Int64)
		if err != nil {
			return "", err
		}
		if !ok {
			return types.FailStatus, nil
		}
	}
	return types.SuccessStatus, nil
}

func (v *validatingStateHandler) verifyCard(ctx context.Context, tx dbclient.TX, cardID, walletID int64) (bool, error) {
	var walletOwnerID int64
	row := tx.QueryRowContext(ctx, selectCard, cardID)
	if err := row.Scan(&walletOwnerID); err != nil {
		return false, err
	}
	return walletID == walletOwnerID, nil
}