package wallet

import (
	"context"
	"errors"
	"ewallet/modules/share"
)

const (
	selectAsset = "SELECT wallet_id, currency, amount FROM wallet_assets WHERE wallet_id = ?"
	selectAssetByCurrency = "SELECT wallet_id, currency, amount FROM wallet_assets WHERE wallet_id = ? AND currency = ?"
)

type (
	AssetDAO struct {
		WalletID string `db:"wallet_id"`
		Amount string `db:"amount"`
		Currency string `db:"currency"`
	}

	Repository interface {
		GetBalance(ctx context.Context, walletID string) ([]AssetDAO, error)
		GetBalanceByCurrency(ctx context.Context, walletID, currency string) (AssetDAO, error)
	}

	repoImpl struct {
		db share.DB
	}
)

func NewRepository(db share.DB) Repository {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) GetBalance(ctx context.Context, walletID string) ([]AssetDAO, error) {
	var assets []AssetDAO
	err := r.db.SelectContext(ctx, &assets, selectAsset, walletID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
func (r repoImpl) GetBalanceByCurrency(ctx context.Context, walletID, currency string) (AssetDAO, error) {
	var assets []AssetDAO
	err := r.db.SelectContext(ctx, &assets, selectAssetByCurrency, walletID, currency)
	if err != nil {
		return AssetDAO{}, err
	}
	if len(assets) == 0 {
		return AssetDAO{}, errors.New("not found")
	}
	return assets[0], nil
}
