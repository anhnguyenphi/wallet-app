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

func TestValidateStateHandler_Handle(t *testing.T) {
	testCases := map[string]struct {
		txFunc         func() (dbclient.TX, sqlmock.Sqlmock)
		transDAO       *types.TransDAO
		ExpectedStatus types.Status
		Err            error
	}{
		"success - withdraw": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id FROM cards").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(1))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				Type:         types.WithdrawType,
				ToCardID:     sql.NullInt64{Int64: 1},
				FromWalletID: sql.NullInt64{Int64: 1},
				Currency:     "USD",
				Amount:       "100",
			},
			ExpectedStatus: types.SuccessStatus,
		},
		"success - deposit": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id FROM cards").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(1))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				Type:         types.WithdrawType,
				FromWalletID: sql.NullInt64{Int64: 1},
				ToCardID:     sql.NullInt64{Int64: 1},
				Currency:     "USD",
				Amount:       "100",
			},
			ExpectedStatus: types.SuccessStatus,
		},
		"invalid card id": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id FROM cards").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(2))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				Type:         types.WithdrawType,
				FromWalletID: sql.NullInt64{Int64: 1},
				ToCardID:     sql.NullInt64{Int64: 1},
				Currency:     "USD",
				Amount:       "100",
			},
			ExpectedStatus: types.FailStatus,
		},
		"failure": {
			txFunc: func() (dbclient.TX, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT wallet_id FROM cards").
					WithArgs(int64(1)).
					WillReturnError(fmt.Errorf("some error"))
				tx, _ := db.Begin()
				return tx, m
			},
			transDAO: &types.TransDAO{
				Type:         types.WithdrawType,
				FromWalletID: sql.NullInt64{Int64: 1},
				ToCardID:     sql.NullInt64{Int64: 1},
				Currency:     "USD",
				Amount:       "100",
			},
			ExpectedStatus: "",
			Err: fmt.Errorf("some error"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tx, m := testCase.txFunc()
			handler := &validatingStateHandler{}
			state, err := handler.Handle(context.TODO(), tx, testCase.transDAO)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedStatus, state)
			if err := m.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
