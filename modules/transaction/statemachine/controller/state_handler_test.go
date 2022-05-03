package controller

import (
	"context"
	"errors"
	"ewallet/modules/share/dbclient"
	dbmock "ewallet/modules/share/dbclient/mocks"
	"ewallet/modules/transaction/statemachine/mocks"
	"ewallet/modules/transaction/statemachine/types"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStateMachineHandler_Handle(t *testing.T) {
	columns := []string{"id", "type", "from_wallet_id", "from_card_id", "to_wallet_id",
		"to_card_id", "external_tx_id", "amount", "currency", "state", "created_at", "updated_at"}
	testCases := map[string]struct{
		dbFunc              func() (dbclient.DB, sqlmock.Sqlmock)
		handlerGetterFunc   func() types.StateHandlerGetter
		stateControllerFunc func() types.StateController
		status              types.Status
		Err                 error
	}{
		"success": {
			dbFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT id, type, from_wallet_id, from_card_id, to_wallet_id, to_card_id, external_tx_id, amount, currency, state, created_at, updated_at from transactions").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, types.TransferType, 1, nil, 2, nil, nil, "100", "USD", types.WithDrawFromSenderState, "today", "today"))
				m.ExpectExec("UPDATE transactions").
					WithArgs(types.DepositToReceiverState, "", int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
				return db, m
			},
			handlerGetterFunc: func() types.StateHandlerGetter {
				m := &mocks.StateHandlerGetter{}
				sh := &mocks.StateHandler{}
				tx := &dbmock.TX{}
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
				sh.On("Begin", mock.Anything).Return(tx, nil)
				sh.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(types.SuccessStatus, nil)
				m.On("GetHandler", types.WithDrawFromSenderState).Return(sh, nil)
				return m
			},
			stateControllerFunc: func() types.StateController {
				m := &mocks.StateController{}
				m.On("NextState", types.WithDrawFromSenderState, types.SuccessStatus).
					Return(types.DepositToReceiverState)
				return m
			},
			status: types.SuccessStatus,
		},
		"handler failed": {
			dbFunc: func() (dbclient.DB, sqlmock.Sqlmock) {
				db, m, _ := sqlmock.New()
				m.ExpectBegin()
				m.ExpectQuery("SELECT id, type, from_wallet_id, from_card_id, to_wallet_id, to_card_id, external_tx_id, amount, currency, state, created_at, updated_at from transactions").
					WithArgs(int64(1)).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(1, types.TransferType, 1, nil, 2, nil, nil, "100", "USD", types.WithDrawFromSenderState, "today", "today"))
				m.ExpectRollback()
				return db, m
			},
			handlerGetterFunc: func() types.StateHandlerGetter {
				m := &mocks.StateHandlerGetter{}
				sh := &mocks.StateHandler{}
				tx := &dbmock.TX{}
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
				sh.On("Begin", mock.Anything).Return(tx, nil)
				sh.On("Handle", mock.Anything, mock.Anything, mock.Anything).
					Return(types.Status(""), errors.New("fail"))
				m.On("GetHandler", types.WithDrawFromSenderState).Return(sh, nil)
				return m
			},
			stateControllerFunc: func() types.StateController {
				m := &mocks.StateController{}
				m.On("NextState", types.WithDrawFromSenderState, types.SuccessStatus).
					Return(types.DepositToReceiverState)
				return m
			},
			status: types.Status(""),
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			db, m := testCase.dbFunc()
			handler := NewStateMachineHandler(db, testCase.handlerGetterFunc(), testCase.stateControllerFunc())
			status, err := handler.Handle(context.TODO(), 1)
			assert.Equal(t,testCase.status, status)
			assert.Equal(t,testCase.Err, err)
			if err := m.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
