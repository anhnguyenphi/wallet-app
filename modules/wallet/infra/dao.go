package infra

import (
	"database/sql"
)

type (
	AssetDAO struct {
		WalletID int64 `db:"wallet_id"`
		Amount string `db:"amount"`
		Currency string `db:"currency"`
	}
)

func RowsToAssetDAOs(rows *sql.Rows) ([]AssetDAO, error) {
	var result []AssetDAO
	for rows.Next() {
		var assetDao AssetDAO
		err := rows.Scan(&assetDao.WalletID, &assetDao.Currency, &assetDao.Amount)
		if err != nil {
			return nil, err
		}
		result = append(result, assetDao)
	}
	return result, nil
}
