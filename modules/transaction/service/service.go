package service

import (
	"context"
	"database/sql"
	"ewallet/modules/asset"
	"ewallet/modules/transaction/infra"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	Service interface {
		Deposit(ctx context.Context, fromCardID, toWalletID int64, amount asset.Asset) (int64, error)
		Withdraw(ctx context.Context, fromWalletID, toCardID int64, amount asset.Asset) (int64, error)
		Transfer(ctx context.Context, fromWalletID, toWalletID int64, amount asset.Asset) (int64, error)
	}

	Publisher interface {
		Publish(ctx context.Context, transType types.Type, transId int64) error
	}

	serviceImpl struct {
		transRepo      infra.Repository
		transPublisher Publisher
	}
)

func NewService(transRepo infra.Repository, transPublisher Publisher) Service {
	return &serviceImpl{
		transRepo: transRepo,
		transPublisher: transPublisher,
	}
}

func (s *serviceImpl) Withdraw(ctx context.Context, fromWalletID, toCardID int64, amount asset.Asset) (int64, error) {
	transID, err := s.transRepo.Create(ctx, types.TransDAO{
		FromWalletID: sql.NullInt64{Int64: fromWalletID},
		ToCardID:     sql.NullInt64{Int64: toCardID},
		Amount:       amount.Amount,
		Currency:     amount.Currency,
		Type:         types.WithdrawType,
	})
	if err != nil {
		return 0, err
	}

	if err := s.transPublisher.Publish(ctx, types.WithdrawType, transID); err != nil {
		return 0, err
	}
	return transID, nil
}

func (s *serviceImpl) Deposit(ctx context.Context, fromCardID, toWalletID int64, amount asset.Asset) (int64, error) {
	transID, err := s.transRepo.Create(ctx, types.TransDAO{
		FromCardID: sql.NullInt64{Int64: fromCardID},
		ToWalletID: sql.NullInt64{Int64: toWalletID},
		Amount:     amount.Amount,
		Currency:   amount.Currency,
		Type:       types.DepositType,
	})
	if err != nil {
		return 0, err
	}

	if err := s.transPublisher.Publish(ctx, types.DepositType, transID); err != nil {
		return 0, err
	}
	return transID, nil
}

func (s *serviceImpl) Transfer(ctx context.Context, fromWalletID, toWalletID int64, amount asset.Asset) (int64, error) {
	transID, err := s.transRepo.Create(ctx, types.TransDAO{
		FromWalletID: sql.NullInt64{Int64: fromWalletID},
		ToWalletID:   sql.NullInt64{Int64: toWalletID},
		Amount:       amount.Amount,
		Currency:     amount.Currency,
		Type:         types.TransferType,
	})
	if err != nil {
		return 0, err
	}

	if err := s.transPublisher.Publish(ctx, types.TransferType, transID); err != nil {
		return 0, err
	}
	return transID, nil
}
