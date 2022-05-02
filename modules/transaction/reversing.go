package transaction

import (
	"context"
	"database/sql"
	"ewallet/modules/share"
)

type (
	reversingStateHandler struct {
		assetDB share.DB
	}
)

func (w *reversingStateHandler) Begin(ctx context.Context) (share.TX, error) {
	return w.assetDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (d *reversingStateHandler) Handle(ctx context.Context,tx share.TX, trans TransDAO) (State, error) {
	var assetDAO AssetDAO
	res := tx.QueryRowContext(ctx, selectWalletForUpdate, trans.FromWalletID, trans.Currency)
	if err := res.Scan(&assetDAO); err != nil {
		return ReversingToSenderState, err
	}
	newBalance, err := add(assetDAO.Amount, trans.Amount)
	if err != nil {
		return ReversingToSenderState, err
	}
	_, err = tx.ExecContext(ctx, updateAmountOfWallet, newBalance, trans.FromWalletID, trans.Currency)
	if err != nil {
		return ReversingToSenderState, err
	}
	return RejectedState, nil
}
