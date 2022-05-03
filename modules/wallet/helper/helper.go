package helper

import "github.com/shopspring/decimal"

func GreaterThanOrEqual(balance, amount string) (bool, error) {
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

