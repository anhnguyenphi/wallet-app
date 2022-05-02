package transaction

import (
	"github.com/shopspring/decimal"
)

func add(balance, amount string) (string, error) {
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

func sub(balance, amount string) (string, error) {
	balanceDecimal, err := decimal.NewFromString(balance)
	if err != nil {
		return "", err
	}
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return "", err
	}
	if balanceDecimal.LessThan(amountDecimal) {
		return "", InsufficientBalance
	}

	return balanceDecimal.Sub(amountDecimal).String(), nil
}

func greaterThanOrEqual(balance, amount string) (bool, error) {
	balanceDecimal, err := decimal.NewFromString(balance)
	if err != nil {
		return false, err
	}
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return false, err
	}
	return balanceDecimal.GreaterThanOrEqual(amountDecimal), nil
}

