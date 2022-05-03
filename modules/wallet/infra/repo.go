package infra

import (
	"context"
	"errors"
	"ewallet/modules/share/dbclient"
)

const (
	selectAsset = "SELECT wallet_id, currency, amount FROM wallet_assets WHERE wallet_id = ?"
	selectAssetByCurrency = "SELECT wallet_id, currency, amount FROM wallet_assets WHERE wallet_id = ? AND currency = ?"
)

type (
	Repository interface {
		GetBalance(ctx context.Context, walletID int64) ([]AssetDAO, error)
		GetBalanceByCurrency(ctx context.Context, walletID int64, currency string) (AssetDAO, error)
	}

	repoImpl struct {
		db dbclient.DB
	}
)

func NewRepository(db dbclient.DB) Repository {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) GetBalance(ctx context.Context, walletID int64) ([]AssetDAO, error) {
	var assets []AssetDAO
	rows, err := r.db.QueryContext(ctx, selectAsset, walletID)
	if err != nil {
		return nil, err
	}
	assets, err = RowsToAssetDAOs(rows)
	if err != nil {
		return nil, err
	}
	return assets, nil
}
func (r repoImpl) GetBalanceByCurrency(ctx context.Context, walletID int64, currency string) (AssetDAO, error) {
	var assets []AssetDAO
	rows, err := r.db.QueryContext(ctx, selectAssetByCurrency, walletID, currency)
	if err != nil {
		return AssetDAO{}, err
	}
	assets, err = RowsToAssetDAOs(rows)
	if err != nil {
		return AssetDAO{}, err
	}
	if len(assets) == 0 {
		return AssetDAO{}, errors.New("not found")
	}
	return assets[0], nil
}
