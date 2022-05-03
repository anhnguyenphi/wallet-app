package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGreaterThanOrEqual(t *testing.T)  {
	testCases := map[string]struct{
		Balance string
		Amount string
		Expected bool
		hasErr bool
	}{
		"greater": {
			Balance: "2",
			Amount: "1",
			Expected: true,
			hasErr: false,
		},
		"equal": {
			Balance: "1",
			Amount: "1",
			Expected: true,
			hasErr: false,
		},
		"less": {
			Balance: "1",
			Amount: "2",
			Expected: false,
			hasErr: false,
		},
		"invalid balance": {
			Balance: "abc",
			Amount: "0",
			Expected: false,
			hasErr: true,
		},
		"invalid amount": {
			Balance: "abc",
			Amount: "0",
			Expected: false,
			hasErr: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := GreaterThanOrEqual(testCase.Balance, testCase.Amount)
			if testCase.hasErr {
				assert.NotNil(t, err)
			}
			assert.Equal(t, testCase.Expected, result)
		})
	}
}
