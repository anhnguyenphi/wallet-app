package state

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/helper"
	"ewallet/modules/transaction/statemachine/types"
)
type (
	withdrawFromWalletStateHandler struct {
		assetDB dbclient.DB
	}
)

func (w *withdrawFromWalletStateHandler) Begin(ctx context.Context) (dbclient.TX, error) {
	return w.assetDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (w *withdrawFromWalletStateHandler) Handle(ctx context.Context,tx dbclient.TX, trans *types.TransDAO) (types.Status, error) {
	var assetDAO types.AssetDAO
	res := tx.QueryRowContext(ctx, selectWalletForUpdate, trans.FromWalletID.Int64, trans.Currency)
	if err := res.Scan(&assetDAO.WalletID, &assetDAO.Amount); err != nil {
		return "", err
	}

	newBalance, err := helper.Sub(assetDAO.Amount, trans.Amount)
	if err != nil {
		return "", err
	}
	negative, err := helper.IsNegative(newBalance)
	if err != nil {
		return "", err
	}
	if negative {
		return types.FailStatus, nil
	}
	_, err = tx.ExecContext(ctx, updateAmountOfWallet, newBalance, trans.FromWalletID.Int64, trans.Currency)
	if err != nil {
		return "", err
	}
	return types.SuccessStatus, nil
}
