package wallet

import (
	"context"
	"errors"
	"ewallet/modules/asset"
)

var (
	ErrEmptyWalletID = errors.New("empty wallet id")
)

type (
	Service interface {
		Deposit(ctx context.Context, sender, receiver Wallet, amount asset.Asset) (int64, error)
		Withdraw(ctx context.Context, sender, receiver Wallet, amount asset.Asset) (int64, error)
		Send(ctx context.Context, sender, receiver Wallet, amount asset.Asset) (int64, error)
		BalanceOf(ctx context.Context, wallet Wallet) ([]asset.Asset, error)
	}

	TransactionService interface {
		Transfer(ctx context.Context, fromWalletID, toWalletID string, amount asset.Asset) (int64, error)
	}

	serviceImpl struct {
		repo               Repository
		transactionService TransactionService
	}
)

func NewService(repo Repository, service TransactionService) Service {
	return &serviceImpl{
		repo: repo,
		transactionService: service,
	}
}

func (s *serviceImpl) Deposit(ctx context.Context,sender, receiver Wallet, amount asset.Asset) (int64, error) {
	panic("implement me")
}

func (s *serviceImpl) Withdraw(ctx context.Context,sender, receiver Wallet, amount asset.Asset) (int64, error) {
	panic("implement me")
}

func (s *serviceImpl) Send(ctx context.Context, sender, receiver Wallet, amount asset.Asset) (int64, error) {
	return s.transactionService.Transfer(ctx, sender.ID, receiver.ID, amount)
}

func (s *serviceImpl) BalanceOf(ctx context.Context,wallet Wallet) ([]asset.Asset, error) {
	if len(wallet.ID) == 0 {
		return nil, ErrEmptyWalletID
	}
	res, err := s.repo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, err
	}
	assets, err := fromAssetDAOToMoney(res)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func fromAssetDAOToMoney(assetDAOs []AssetDAO) ([]asset.Asset, error) {
	assets := make([]asset.Asset, len(assetDAOs))
	for idx, assetDAO := range assetDAOs {
		assets[idx] = asset.Asset{
			WalletID: assetDAO.WalletID,
			Amount: assetDAO.Amount,
			Currency: assetDAO.Currency,
		}
	}
	return assets, nil
}
