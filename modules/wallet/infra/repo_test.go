package infra

import (
	"context"
	"errors"
	"ewallet/modules/share/dbclient"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoImpl_GetBalance(t *testing.T) {
	columns := []string{"wallet_id", "currency", "amount"}
	testCases := map[string]struct{
		DBFunc func() (dbclient.DB, sqlmock.Sqlmock)
		WalletID int64
		ExpectedAsset []AssetDAO
		Err error
	}{
		"success": {
			DBFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				dbMock, m, _ := sqlmock.New()
				m.ExpectQuery("SELECT wallet_id, currency, amount FROM wallet_assets").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "USD", "100"))
				return dbMock, m
			},
			WalletID: 1,
			ExpectedAsset: []AssetDAO{
				{
					WalletID: 1,
					Currency: "USD",
					Amount:   "100",
				},
			},
			Err: nil,
		},
		"fail": {
			DBFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				dbMock, m, _ := sqlmock.New()
				m.ExpectQuery("SELECT wallet_id, currency, amount FROM wallet_assets").
					WithArgs(int64(1)).
					WillReturnError(errors.New("fail"))
				return dbMock, m
			},
			WalletID: 1,
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db, m := testCase.DBFunc()
			repo := NewRepository(db)
			result, err := repo.GetBalance(context.TODO(), testCase.WalletID)
			assert.Equal(t,testCase.Err, err)
			assert.Equal(t,testCase.ExpectedAsset, result)
			if err := m.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepoImpl_GetBalanceByCurrency(t *testing.T) {
	columns := []string{"wallet_id", "currency", "amount"}
	testCases := map[string]struct{
		DBFunc func() (dbclient.DB, sqlmock.Sqlmock)
		WalletID      int64
		Currency      string
		ExpectedAsset AssetDAO
		Err           error
	}{
		"success": {
			DBFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				dbMock, m, _ := sqlmock.New()
				m.ExpectQuery("SELECT wallet_id, currency, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "USD", "100"))
				return dbMock, m
			},
			WalletID: 1,
			Currency: "USD",
			ExpectedAsset: AssetDAO{
				WalletID: 1,
				Currency: "USD",
				Amount:   "100",
			},
			Err: nil,
		},
		"fail": {
			DBFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				dbMock, m, _ := sqlmock.New()
				m.ExpectQuery("SELECT wallet_id, currency, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnError(errors.New("fail"))
				return dbMock, m
			},
			Currency: "USD",
			WalletID: 1,
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db, m := testCase.DBFunc()
			repo := NewRepository(db)
			result, err := repo.GetBalanceByCurrency(context.TODO(), testCase.WalletID, testCase.Currency)
			assert.Equal(t,testCase.Err, err)
			assert.Equal(t,testCase.ExpectedAsset, result)
			if err := m.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
