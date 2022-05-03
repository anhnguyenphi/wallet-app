package state

import (
	"context"
	"database/sql"
	"errors"
	dbmock "ewallet/modules/share/dbclient/mocks"
	"ewallet/modules/transaction/statemachine/mocks"
	"ewallet/modules/transaction/statemachine/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDepositToBankStateHandler_Handle(t *testing.T) {
	testCases := map[string]struct{
		bankDepositFunc func() BankDeposit
		transDAO        *types.TransDAO
		ExpectedStatus  types.Status
		ExTransactionID string
		Err             error
	}{
		"success": {
			bankDepositFunc: func() BankDeposit {
				bd := &mocks.BankDeposit{}
				bd.On("CreateDepositTransaction", mock.Anything, int64(1), "100").
					Return("abc", nil)
				return bd
			},
			transDAO: &types.TransDAO{
				FromCardID: sql.NullInt64{Int64: 1},
				Amount: "100",
			},
			ExpectedStatus:  types.SuccessStatus,
			ExTransactionID: "abc",
		},
		"failure": {
			bankDepositFunc: func() BankDeposit {
				bd := &mocks.BankDeposit{}
				bd.On("CreateDepositTransaction", mock.Anything, int64(1), "100").
					Return("", errors.New("fail"))
				return bd
			},
			transDAO: &types.TransDAO{
				FromCardID: sql.NullInt64{Int64: 1},
				Amount: "100",
			},
			ExpectedStatus:  types.FailStatus,
			ExTransactionID: "",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			handler := &depositToBankStateHandler{bank: testCase.bankDepositFunc()}
			state, err := handler.Handle(context.TODO(), &dbmock.TX{}, testCase.transDAO)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedStatus, state)
			assert.Equal(t, testCase.ExTransactionID, testCase.transDAO.ExTransactionID.String)

		})
	}
}
