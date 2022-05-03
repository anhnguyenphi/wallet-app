package service

import (
	"context"
	"errors"
	"ewallet/modules/asset"
	"ewallet/modules/wallet/infra"
	"ewallet/modules/wallet/infra/mocks"
	serviceMock "ewallet/modules/wallet/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestServiceImpl_BalanceOf(t *testing.T) {
	testCases := map[string]struct{
		RepoFunc func() infra.Repository
		WalletID int64
		ExpectedAsset []asset.Asset
		Err error
	}{
		"success": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalance", mock.Anything, int64(1)).
					Return([]infra.AssetDAO{
					{
						WalletID: 1,
						Currency: "USD",
						Amount:   "100",
					},
				}, nil)
				return repoMock
			},
			WalletID: 1,
			ExpectedAsset: []asset.Asset{
				{
					WalletID: 1,
					Currency: "USD",
					Amount:   "100",
				},
			},
			Err: nil,
		},
		"fail": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalance", mock.Anything, int64(1)).
					Return(nil, errors.New("fail"))
				return repoMock
			},
			WalletID: 1,
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			service := NewService(testCase.RepoFunc(), nil)
			result, err := service.BalanceOf(context.TODO(), testCase.WalletID)
			assert.Equal(t,testCase.Err, err)
			assert.Equal(t,testCase.ExpectedAsset, result)
		})
	}
}

func TestServiceImpl_Deposit(t *testing.T) {
	testCases := map[string]struct{
		TransactionServiceFunc func() TransactionService
		FromCardID int64
		ToWalletID int64
		Amount asset.Asset
		ExpectedTransID int64
		Err error
	}{
		"success": {
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Deposit", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(1), nil)
				return ts
			},
			FromCardID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			ExpectedTransID: int64(1),
		},
		"fail": {
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Deposit", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(0), errors.New("fail"))
				return ts
			},
			FromCardID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			Err: errors.New("fail"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			service := NewService(nil, testCase.TransactionServiceFunc())
			result, err := service.Deposit(context.TODO(), testCase.FromCardID, testCase.ToWalletID, testCase.Amount)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, result)
		})
	}
}

func TestServiceImpl_Withdraw(t *testing.T) {
	testCases := map[string]struct{
		RepoFunc func() infra.Repository
		TransactionServiceFunc func() TransactionService
		FromWalletID int64
		ToCardID int64
		Amount asset.Asset
		ExpectedTransID int64
		Err error
	}{
		"success": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "200",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Withdraw", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(1), nil)
				return ts
			},
			FromWalletID: int64(1),
			ToCardID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			ExpectedTransID: int64(1),
		},
		"insufficient balance": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "80",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Withdraw", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(1), nil)
				return ts
			},
			FromWalletID: int64(1),
			ToCardID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "200",
			},
			ExpectedTransID: int64(0),
			Err: InsufficientBalance,
		},
		"can not get balance": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{}, errors.New("repo failed"))
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Withdraw", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(0), errors.New("fail"))
				return ts
			},
			FromWalletID: int64(1),
			ToCardID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			Err: errors.New("repo failed"),
		},
		"fail to withdraw": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "200",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Withdraw", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(0), errors.New("fail to withdraw"))
				return ts
			},
			FromWalletID: int64(1),
			ToCardID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			Err: errors.New("fail to withdraw"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			service := NewService(testCase.RepoFunc(), testCase.TransactionServiceFunc())
			result, err := service.Withdraw(context.TODO(), testCase.FromWalletID, testCase.ToCardID, testCase.Amount)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, result)
		})
	}
}

func TestServiceImpl_Transfer(t *testing.T) {
	testCases := map[string]struct{
		RepoFunc func() infra.Repository
		TransactionServiceFunc func() TransactionService
		FromWalletID int64
		ToWalletID int64
		Amount asset.Asset
		ExpectedTransID int64
		Err error
	}{
		"success": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "200",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Transfer", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(1), nil)
				return ts
			},
			FromWalletID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			ExpectedTransID: int64(1),
		},
		"insufficient balance": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "80",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Transfer", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(1), nil)
				return ts
			},
			FromWalletID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "200",
			},
			ExpectedTransID: int64(0),
			Err: InsufficientBalance,
		},
		"can not get balance": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{}, errors.New("repo failed"))
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Transfer", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(0), errors.New("fail"))
				return ts
			},
			FromWalletID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			Err: errors.New("repo failed"),
		},
		"fail to withdraw": {
			RepoFunc: func() infra.Repository {
				repoMock := &mocks.Repository{}
				repoMock.On("GetBalanceByCurrency", mock.Anything, int64(1), "USD").
					Return(infra.AssetDAO{
						WalletID: 1,
						Currency: "USD",
						Amount:   "200",
					}, nil)
				return repoMock
			},
			TransactionServiceFunc: func() TransactionService {
				ts := &serviceMock.TransactionService{}
				ts.On("Transfer", mock.Anything, int64(1), int64(2), asset.Asset{
					Currency: "USD",
					Amount: "100",
				}).Return(int64(0), errors.New("fail to withdraw"))
				return ts
			},
			FromWalletID: int64(1),
			ToWalletID: int64(2),
			Amount: asset.Asset{
				Currency: "USD",
				Amount: "100",
			},
			Err: errors.New("fail to withdraw"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			service := NewService(testCase.RepoFunc(), testCase.TransactionServiceFunc())
			result, err := service.Transfer(context.TODO(), testCase.FromWalletID, testCase.ToWalletID, testCase.Amount)
			assert.Equal(t, testCase.Err, err)
			assert.Equal(t, testCase.ExpectedTransID, result)
		})
	}
}