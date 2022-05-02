package transaction

import (
	"context"
	"errors"
	"ewallet/modules/asset"
	"ewallet/modules/wallet"
)

var (
	InsufficientBalance = errors.New("insufficient balance")
)

type (
	Status string

	Service interface {
		Transfer(ctx context.Context, fromWalletID, toWalletID string, amount asset.Asset) (int64, error)
	}

	WalletRepository interface {
		GetBalance(ctx context.Context, walletID string) ([]wallet.AssetDAO, error)
		GetBalanceByCurrency(ctx context.Context, walletID, currency string) (wallet.AssetDAO, error)
	}

	serviceImpl struct {
		transRepo Repository
		walletRepo WalletRepository
	}
)

func NewService(transRepo Repository, walletRepo WalletRepository) Service {
	return &serviceImpl{
		transRepo: transRepo,
		walletRepo: walletRepo,
	}
}

func (s *serviceImpl) Transfer(ctx context.Context, fromWalletID, toWalletID string, amount asset.Asset) (int64, error) {
	balance, err := s.walletRepo.GetBalanceByCurrency(ctx, fromWalletID, amount.Currency)
	if err != nil {
		return 0, err
	}
	// normal check to make sure we don't create too many failed transaction
	if isOk, err := greaterThanOrEqual(balance.Amount, amount.Amount); err != nil || !isOk {
		return 0, InsufficientBalance
	}
	return s.transRepo.Create(ctx, TransDAO{
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		Amount:       amount.Amount,
		Currency:     amount.Currency,
	})
}
