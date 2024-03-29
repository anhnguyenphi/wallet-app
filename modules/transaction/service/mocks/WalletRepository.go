// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	infra "ewallet/modules/wallet/infra"

	mock "github.com/stretchr/testify/mock"
)

// WalletRepository is an autogenerated mock type for the WalletRepository type
type WalletRepository struct {
	mock.Mock
}

// GetBalance provides a mock function with given fields: ctx, walletID
func (_m *WalletRepository) GetBalance(ctx context.Context, walletID int64) ([]infra.AssetDAO, error) {
	ret := _m.Called(ctx, walletID)

	var r0 []infra.AssetDAO
	if rf, ok := ret.Get(0).(func(context.Context, int64) []infra.AssetDAO); ok {
		r0 = rf(ctx, walletID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]infra.AssetDAO)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, walletID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalanceByCurrency provides a mock function with given fields: ctx, walletID, currency
func (_m *WalletRepository) GetBalanceByCurrency(ctx context.Context, walletID int64, currency string) (infra.AssetDAO, error) {
	ret := _m.Called(ctx, walletID, currency)

	var r0 infra.AssetDAO
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) infra.AssetDAO); ok {
		r0 = rf(ctx, walletID, currency)
	} else {
		r0 = ret.Get(0).(infra.AssetDAO)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, walletID, currency)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
