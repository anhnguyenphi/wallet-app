// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// BankTransactionVerifier is an autogenerated mock type for the BankTransactionVerifier type
type BankTransactionVerifier struct {
	mock.Mock
}

// GetStatus provides a mock function with given fields: ctx, cardID, bankingTransactionID
func (_m *BankTransactionVerifier) GetStatus(ctx context.Context, cardID int64, bankingTransactionID string) (string, error) {
	ret := _m.Called(ctx, cardID, bankingTransactionID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) string); ok {
		r0 = rf(ctx, cardID, bankingTransactionID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, cardID, bankingTransactionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
