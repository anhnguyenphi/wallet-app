package types

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
)

type (
	Status string
	State string
	Type string

	StateController interface {
		NextState(state State, status Status) State
	}

	StateHandler interface {
		Handle(ctx context.Context,tx dbclient.TX, trans *TransDAO) (Status, error)
		Begin(ctx context.Context) (dbclient.TX, error)
	}

	StateHandlerGetter interface {
		GetHandler(state State) (StateHandler, error)
	}

	AssetDAO struct {
		WalletID string `db:"wallet_id"`
		Amount string `db:"amount"`
		Currency string `db:"currency"`
	}

	TransDAO struct {
		ID              string         `db:"id"`
		Type            Type           `db:"type"`
		FromWalletID    sql.NullInt64  `db:"from_wallet_id"`
		FromCardID      sql.NullInt64  `db:"from_card_id"`
		ToWalletID      sql.NullInt64  `db:"to_wallet_id"`
		ToCardID        sql.NullInt64  `db:"to_card_id"`
		ExTransactionID sql.NullString `db:"external_tx_id"`
		Amount          string         `db:"amount"`
		Currency        string         `db:"currency"`
		State           State          `db:"state"`
		CreatedAt       string         `db:"created_at"`
		UpdatedAt       string         `db:"updated_at"`
	}
)

func (trans *TransDAO) Scan(row *sql.Row) error {
	return row.Scan(
		&trans.ID,
		&trans.Type,
		&trans.FromWalletID,
		&trans.FromCardID,
		&trans.ToWalletID,
		&trans.ToCardID,
		&trans.ExTransactionID,
		&trans.Amount,
		&trans.Currency,
		&trans.State,
		&trans.CreatedAt,
		&trans.UpdatedAt)
}