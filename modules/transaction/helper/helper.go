package helper

import (
	"github.com/shopspring/decimal"
)

func Add(balance, amount string) (string, error) {
	balanceDecimal, err := decimal.NewFromString(balance)
	if err != nil {
		return "", err
	}
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return "", err
	}

	return balanceDecimal.Add(amountDecimal).String(), nil
}

func Sub(balance, amount string) (string, error) {
	balanceDecimal, err := decimal.NewFromString(balance)
	if err != nil {
		return "", err
	}
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return "", err
	}

	return balanceDecimal.Sub(amountDecimal).String(), nil
}

func IsNegative(balance string) (bool, error) {
	balanceDecimal, err := decimal.NewFromString(balance)
	if err != nil {
		return false, err
	}
	zero := decimal.NewFromInt(0)
	return balanceDecimal.LessThan(zero), nil
}
