package service

import (
	"context"
	"encoding/hex"
	"math/rand"
)

const (
	defaultTrxLen = 10
	StatusSuccess = "success"
	StatusInProgress = "in_progress"
	StatusFailed = "failed"
)

type (
	Bank interface {

	}

	bank struct {

	}
)

func NewBank() *bank {
	return &bank{}
}

func (b *bank) CreateDepositTransaction(ctx context.Context, cardID int64, amount string) (string, error) {
	return generateID(defaultTrxLen), nil
}

func (b *bank) CreateWithdrawTransaction(ctx context.Context, cardID int64, amount string) (string, error) {
	return generateID(defaultTrxLen), nil
}

func (b *bank) GetStatus(ctx context.Context, cardID int64, bankingTransactionID string) (string, error) {
	return StatusSuccess, nil
}

func generateID(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

