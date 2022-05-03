package service

import (
	"context"
	"errors"
	"ewallet/modules/asset"
	"ewallet/modules/wallet/helper"
	"ewallet/modules/wallet/infra"
)

var (
	InsufficientBalance = errors.New("insufficient balance")
)

type (
	Service interface {
		Deposit(ctx context.Context, fromCardID, toWalletID int64, amount asset.Asset) (int64, error)
		Withdraw(ctx context.Context, fromWalletID, toCardID int64, amount asset.Asset) (int64, error)
		Transfer(ctx context.Context, fromWalletID, toWalletID int64, amount asset.Asset) (int64, error)
		BalanceOf(ctx context.Context, walletID int64) ([]asset.Asset, error)
	}

	TransactionService interface {
		Deposit(ctx context.Context, fromCardID, toWalletID int64, amount asset.Asset) (int64, error)
		Withdraw(ctx context.Context, fromWalletID, toCardID int64, amount asset.Asset) (int64, error)
		Transfer(ctx context.Context, fromWalletID, toWalletID int64, amount asset.Asset) (int64, error)
	}

	serviceImpl struct {
		repo               infra.Repository
		transactionService TransactionService
	}
)

func NewService(repo infra.Repository, service TransactionService) Service {
	return &serviceImpl{
		repo: repo,
		transactionService: service,
	}
}

func (s *serviceImpl) Deposit(ctx context.Context, fromCardID, toWalletID int64, amount asset.Asset) (int64, error) {
	return s.transactionService.Deposit(ctx, fromCardID, toWalletID, amount)
}

func (s *serviceImpl) Withdraw(ctx context.Context, fromWalletID, toCardID int64, amount asset.Asset) (int64, error) {
	balance, err := s.repo.GetBalanceByCurrency(ctx, fromWalletID, amount.Currency)
	if err != nil {
		return 0, err
	}
	// normal check to make sure we don't create too many failed transaction
	if isOk, err := helper.GreaterThanOrEqual(balance.Amount, amount.Amount); err != nil || !isOk {
		return 0, InsufficientBalance
	}

	return s.transactionService.Withdraw(ctx, fromWalletID, toCardID, amount)
}

func (s *serviceImpl) Transfer(ctx context.Context, fromWalletID, toWalletID int64, amount asset.Asset) (int64, error) {
	balance, err := s.repo.GetBalanceByCurrency(ctx, fromWalletID, amount.Currency)
	if err != nil {
		return 0, err
	}
	// normal check to make sure we don't create too many failed transaction
	if isOk, err := helper.GreaterThanOrEqual(balance.Amount, amount.Amount); err != nil || !isOk {
		return 0, InsufficientBalance
	}

	return s.transactionService.Transfer(ctx, fromWalletID, toWalletID, amount)
}

func (s *serviceImpl) BalanceOf(ctx context.Context, walletID int64) ([]asset.Asset, error) {
	res, err := s.repo.GetBalance(ctx, walletID)
	if err != nil {
		return nil, err
	}
	assets, err := fromAssetDAOToMoney(res)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func fromAssetDAOToMoney(assetDAOs []infra.AssetDAO) ([]asset.Asset, error) {
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
