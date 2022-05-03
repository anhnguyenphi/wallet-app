package state

import (
	"context"
	"database/sql"
	"ewallet/modules/bank/service"
	dbmock "ewallet/modules/share/dbclient/mocks"
	"ewallet/modules/transaction/statemachine/mocks"
	"ewallet/modules/transaction/statemachine/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestVerifyBankingTrxStateHandler_Handle(t *testing.T) {
	testCases := map[string]struct{
		bankVerifierFunc func() BankTransactionVerifier
		transDAO         *types.TransDAO
		ExpectedStatus   types.Status
		Err              error
	}{
		"success - withdraw": {
			bankVerifierFunc: func() BankTransactionVerifier {
				bd := &mocks.BankTransactionVerifier{}
				bd.On("GetStatus", mock.Anything, int64(1), "abc").
					Return(service.StatusSuccess, nil)
				return bd
			},
			transDAO: &types.TransDAO{
				FromCardID: sql.NullInt64{Int64: 1},
				Amount: "100",
				ExTransactionID: sql.NullString{String: "abc"},
			},
			ExpectedStatus: types.SuccessStatus,
		},
		"success - deposit": {
			bankVerifierFunc: func() BankTransactionVerifier {
				bd := &mocks.BankTransactionVerifier{}
				bd.On("GetStatus", mock.Anything, int64(1), "abc").
					Return(service.StatusSuccess, nil)
				return bd
			},
			transDAO: &types.TransDAO{
				FromCardID: sql.NullInt64{Int64: 1},
				Amount: "100",
				ExTransactionID: sql.NullString{String: "abc"},
			},
			ExpectedStatus: types.SuccessStatus,
		},
		"failure": {
			bankVerifierFunc: func() BankTransactionVerifier {
				bd := &mocks.BankTransactionVerifier{}
				bd.On("GetStatus", mock.Anything, int64(1), "abc").
					Return(service.StatusFailed, nil)
				return bd
			},
			transDAO: &types.TransDAO{
				FromCardID: sql.NullInt64{Int64: 1},
				Amount: "100",
				ExTransactionID: sql.NullString{String: "abc"},
			},
			ExpectedStatus: types.FailStatus,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			handler := &verifyBankingTrxStateHandler{bank: testCase.bankVerifierFunc()}
			state, err := handler.Handle(context.TODO(), &dbmock.TX{}, testCase.transDAO)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedStatus, state)

		})
	}
}
