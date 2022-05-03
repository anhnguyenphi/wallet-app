package state

import (
	"context"
	"ewallet/modules/transaction/statemachine/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompleteStateHandler_Handle(t *testing.T) {
	handler := &completeStateHandler{}
	status, err := handler.Handle(context.TODO(), nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, types.CompleteStatus, status)
}
