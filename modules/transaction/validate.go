package transaction

import (
	"context"
	"database/sql"
	"ewallet/modules/share"
)

type (
	validatingStateHandler struct {
		assetDB share.DB
	}
)

func (w *validatingStateHandler) Begin(ctx context.Context) (share.TX, error) {
	return w.assetDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (v *validatingStateHandler) Handle(ctx context.Context,tx share.TX, trans TransDAO) (State, error) {
	var FromAssetDAO, ToAssetDAO AssetDAO
	fromRes := tx.QueryRowContext(ctx, selectWallet, trans.FromWalletID, trans.Currency)
	if err := fromRes.Scan(&FromAssetDAO); err != nil {
		return RejectedState, err
	}
	toRes := tx.QueryRowContext(ctx, selectWallet, trans.FromWalletID, trans.Currency)
	if err := toRes.Scan(&ToAssetDAO); err != nil {
		return RejectedState, err
	}

	return WithDrawFromSenderState, nil
}