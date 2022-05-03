package controller

import (
	"ewallet/modules/transaction/statemachine/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateController_NextState(t *testing.T) {
	controller := &stateController{
		successMap: map[types.State]types.State{
			types.DepositToReceiverState: types.CompletedState,
		},
		failureMap: map[types.State]types.State{
			types.DepositToReceiverState: types.RejectedState,
		},
	}

	testCases := map[string]struct{
		state     types.State
		status    types.Status
		nextState types.State
	}{
		"success": {
			state:     types.DepositToReceiverState,
			status:    types.SuccessStatus,
			nextState: types.CompletedState,
		},
		"fail": {
			state:     types.DepositToReceiverState,
			status:    types.FailStatus,
			nextState: types.RejectedState,
		},
		"complete": {
			state:     types.CompletedState,
			status:    types.CompleteStatus,
			nextState: types.CompletedState,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			state := controller.NextState(testCase.state, testCase.status)
			assert.Equal(t, testCase.nextState, state)
		})
	}
}
