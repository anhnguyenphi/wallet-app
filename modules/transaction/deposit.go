package transaction

import (
	"context"
	"database/sql"
	"ewallet/modules/share"
)

type (
	depositStateHandler struct {
		assetDB share.DB
	}
)

func (w *depositStateHandler) Begin(ctx context.Context) (share.TX, error) {
	return w.assetDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (d *depositStateHandler) Handle(ctx context.Context, tx share.TX, trans TransDAO) (State, error) {
	var assetDAO AssetDAO
	res := tx.QueryRowContext(ctx, selectWalletForUpdate, trans.ToWalletID, trans.Currency)
	if err := res.Scan(&assetDAO); err != nil {
		return DepositToReceiverState, err
	}

	newBalance, err := add(assetDAO.Amount, trans.Amount)
	if err != nil {
		return ReversingToSenderState, err
	}
	_, err = tx.ExecContext(ctx, updateAmountOfWallet, newBalance, trans.ToWalletID, trans.Currency)
	if err != nil {
		return DepositToReceiverState, err
	}
	return CompletedState, nil
}


