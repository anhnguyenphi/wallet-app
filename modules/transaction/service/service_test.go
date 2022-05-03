package service

import (
	"context"
	"database/sql"
	"errors"
	"ewallet/modules/asset"
	"ewallet/modules/transaction/infra"
	infra_mocks "ewallet/modules/transaction/infra/mocks"
	"ewallet/modules/transaction/service/mocks"
	"ewallet/modules/transaction/statemachine/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceImpl_Withdraw(t *testing.T) {
	testCases := map[string]struct{
		repo            func() infra.Repository
		publisherFuc    func() Publisher
		ExpectedTransID int64
		Err             error
	}{
		"success": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.WithdrawType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.WithdrawType, int64(1)).
					Return(nil)
				return m
			},
			ExpectedTransID: int64(1),
		},
		"fail - create tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.WithdrawType,
				}).Return(int64(0), errors.New("fail"))
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.WithdrawType, int64(1)).
					Return(nil)
				return m
			},
			Err: errors.New("fail"),
		},
		"fail - publish tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.WithdrawType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.WithdrawType, int64(1)).
					Return(errors.New("fail"))
				return m
			},
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			serv := NewService(testCase.repo(), testCase.publisherFuc())
			transID, err := serv.Withdraw(context.TODO(), 1, 1, asset.Asset{
				Amount: "100",
				Currency: "USD",
			})
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, transID)
		})
	}
}

func TestServiceImpl_Deposit(t *testing.T) {
	testCases := map[string]struct{
		repo            func() infra.Repository
		publisherFuc    func() Publisher
		ExpectedTransID int64
		Err             error
	}{
		"success": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.DepositType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.DepositType, int64(1)).
					Return(nil)
				return m
			},
			ExpectedTransID: int64(1),
		},
		"fail - create tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.DepositType,
				}).Return(int64(0), errors.New("fail"))
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.DepositType, int64(1)).
					Return(nil)
				return m
			},
			Err: errors.New("fail"),
		},
		"fail - publish tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.DepositType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.DepositType, int64(1)).
					Return(errors.New("fail"))
				return m
			},
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			serv := NewService(testCase.repo(), testCase.publisherFuc())
			transID, err := serv.Deposit(context.TODO(), 1, 1, asset.Asset{
				Amount: "100",
				Currency: "USD",
			})
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, transID)
		})
	}
}

func TestServiceImpl_Transfer(t *testing.T) {
	testCases := map[string]struct{
		repo            func() infra.Repository
		publisherFuc    func() Publisher
		ExpectedTransID int64
		Err             error
	}{
		"success": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.TransferType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.TransferType, int64(1)).
					Return(nil)
				return m
			},
			ExpectedTransID: int64(1),
		},
		"fail - create tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.TransferType,
				}).Return(int64(0), errors.New("fail"))
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.TransferType, int64(1)).
					Return(nil)
				return m
			},
			Err: errors.New("fail"),
		},
		"fail - publish tx": {
			repo: func() infra.Repository {
				m := &infra_mocks.Repository{}
				m.On("Create", mock.Anything, types.TransDAO{
					FromWalletID: sql.NullInt64{Int64: 1},
					ToCardID:     sql.NullInt64{Int64: 1},
					Amount:       "100",
					Currency:     "USD",
					Type:         types.TransferType,
				}).Return(int64(1), nil)
				return m
			},
			publisherFuc: func() Publisher {
				m := &mocks.Publisher{}
				m.On("Publish",mock.Anything, types.TransferType, int64(1)).
					Return(errors.New("fail"))
				return m
			},
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			serv := NewService(testCase.repo(), testCase.publisherFuc())
			transID, err := serv.Transfer(context.TODO(), 1, 1, asset.Asset{
				Amount: "100",
				Currency: "USD",
			})
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, transID)
		})
	}
}