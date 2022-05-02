package transaction

import (
	"context"
	"ewallet/modules/share"
)

const (
	insertTransactionSQL = `INSERT INTO transactions(from_wallet_id, to_wallet_id, amount, currency, status) VALUES(:from_wallet_id, :to_wallet_id, :amount, :currency, :status)`
)

type (

	AssetDAO struct {
		WalletID string `db:"wallet_id"`
		Amount string `db:"amount"`
		Currency string `db:"currency"`
	}

	TransDAO struct {
		ID           string `db:"id"`
		Type         string `db:"type"`
		FromWalletID string `db:"from_wallet_id"`
		ToWalletID   string `db:"to_wallet_id"`
		Amount       string `db:"amount"`
		Currency     string `db:"currency"`
		State        State  `db:"state"`
		CreatedAt    string `db:"created_at"`
		UpdatedAt    string `db:"updated_at"`
	}

	Publisher interface {
		Publish(ctx context.Context, transId int64) error
	}

	Repository interface {
		Create(ctx context.Context, dao TransDAO) (int64, error)
	}

	repo struct {
		db share.DB
		transPublisher Publisher
	}
)

func NewRepository(db share.DB, transPublisher Publisher) Repository {
	return &repo{
		db: db,
		transPublisher: transPublisher,
	}
}

func (r *repo) Create(ctx context.Context, dao TransDAO) (int64, error) {
	tx ,err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	dao.State = ValidatingState
	res, err := tx.ExecContext(ctx, insertTransactionSQL, dao)
	if err != nil {
		return 0, err
	}
	transID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = r.transPublisher.Publish(ctx, transID)
	if err != nil {
		return 0, err
	}

	return transID, tx.Commit()
}
