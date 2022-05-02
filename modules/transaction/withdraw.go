package transaction

import (
	"context"
	"database/sql"
	"ewallet/modules/share"
)
type (
	withdrawStateHandler struct {
		assetDB share.DB
	}
)

func (w *withdrawStateHandler) Begin(ctx context.Context) (share.TX, error) {
	return w.assetDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (w *withdrawStateHandler) Handle(ctx context.Context,tx share.TX, trans TransDAO) (State, error) {
	var assetDAO AssetDAO
	res := tx.QueryRowContext(ctx, selectWalletForUpdate, trans.FromWalletID, trans.Currency)
	if err := res.Scan(&assetDAO); err != nil {
		return WithDrawFromSenderState, err
	}

	newBalance, err := sub(assetDAO.Amount, trans.Amount)
	if err != nil {
		return RejectedState, err
	}
	_, err = tx.ExecContext(ctx, updateAmountOfWallet, newBalance, trans.FromWalletID, trans.Currency)
	if err != nil {
		return WithDrawFromSenderState, err
	}
	return DepositToReceiverState, nil
}
