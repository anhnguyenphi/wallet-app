package state

import (
	"context"
	"database/sql"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithdrawFromWalletStateHandler_Handle(t *testing.T) {
	testCases := map[string]struct {
		txFunc         func() (dbclient.TX, sqlmock.Sqlmock)
		transDAO       *types.TransDAO
		ExpectedStatus types.Status
		Err            error
	}{
		"success": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id", "amount"}).AddRow(1, "100"))
				m.ExpectExec("UPDATE wallet_assets").
					WithArgs("0", int64(1), "USD").
					WillReturnResult(sqlmock.NewResult(1, 1))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				FromWalletID: sql.NullInt64{Int64: 1},
				Currency:   "USD",
				Amount:     "100",
			},
			ExpectedStatus: types.SuccessStatus,
		},
		"failure to lock": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnError(fmt.Errorf("some error"))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				FromWalletID: sql.NullInt64{Int64: 1},
				Currency:   "USD",
				Amount:     "100",
			},
			ExpectedStatus: "",
			Err: fmt.Errorf("some error"),
		},
		"insufficient balance": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id", "amount"}).AddRow(1, "50"))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				FromWalletID: sql.NullInt64{Int64: 1},
				Currency:   "USD",
				Amount:     "100",
			},
			ExpectedStatus: types.FailStatus,
		},
		"failure to update amount": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id, amount FROM wallet_assets").
					WithArgs(int64(1), "USD").
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id", "amount"}).AddRow(1, "100"))
				m.ExpectExec("UPDATE wallet_assets").
					WithArgs("0", int64(1), "USD").
					WillReturnError(fmt.Errorf("some error"))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				FromWalletID: sql.NullInt64{Int64: 1},
				Currency:   "USD",
				Amount:     "100",
			},
			ExpectedStatus: "",
			Err: fmt.Errorf("some error"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tx, m := testCase.txFunc()
			handler := &withdrawFromWalletStateHandler{}
			state, err := handler.Handle(context.TODO(), tx, testCase.transDAO)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedStatus, state)
			if err := m.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
