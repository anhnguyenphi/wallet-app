package infra

import (
	"context"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

const (
	insertTransactionSQL = `INSERT INTO transactions(type, from_wallet_id, from_card_id, to_wallet_id, to_card_id, amount, currency, state) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
)

type (
	Repository interface {
		Create(ctx context.Context, dao types.TransDAO) (int64, error)
	}

	repo struct {
		db dbclient.DB
	}
)

func NewRepository(db dbclient.DB) Repository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, dao types.TransDAO) (int64, error) {
	tx ,err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	dao.State = types.ValidatingState
	res, err := tx.ExecContext(ctx, insertTransactionSQL, dao.Type, dao.FromWalletID.Int64,
		dao.FromCardID.Int64, dao.ToWalletID.Int64, dao.ToCardID.Int64, dao.Amount, dao.Currency, dao.State)
	if err != nil {
		return 0, err
	}
	transID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return transID, tx.Commit()
}
