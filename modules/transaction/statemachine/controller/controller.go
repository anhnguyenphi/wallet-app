package controller

import (
	"ewallet/modules/transaction/statemachine/types"
)

type (
	stateController struct {
		successMap map[types.State]types.State
		failureMap map[types.State]types.State
	}
)

func NewTransferStateController() types.StateController {
	return &stateController{
		successMap: map[types.State]types.State{
			types.ValidatingState:         types.WithDrawFromSenderState,
			types.WithDrawFromSenderState: types.DepositToReceiverState,
			types.DepositToReceiverState:  types.CompletedState,
			types.RefundToSenderState:     types.RejectedState,
		},
		failureMap: map[types.State]types.State{
			types.ValidatingState:         types.RejectedState,
			types.WithDrawFromSenderState: types.RejectedState,
			types.DepositToReceiverState:  types.RefundToSenderState,
			types.RefundToSenderState:     types.RejectedState,
		},
	}
}

func NewDepositStateController() types.StateController {
	return &stateController{
		successMap: map[types.State]types.State{
			types.ValidatingState:               types.WithDrawFromBankState,
			types.WithDrawFromBankState:         types.VerifyBankingTransactionState,
			types.VerifyBankingTransactionState: types.DepositToReceiverState,
			types.DepositToReceiverState:        types.CompletedState,
		},
		failureMap: map[types.State]types.State{
			types.ValidatingState:               types.RejectedState,
			types.WithDrawFromBankState:         types.RejectedState,
			types.VerifyBankingTransactionState: types.RejectedState,
			types.DepositToReceiverState:        types.RejectedState,
		},
	}
}

func NewWithdrawStateController() types.StateController {
	return &stateController{
		successMap: map[types.State]types.State{
			types.ValidatingState:               types.WithDrawFromSenderState,
			types.WithDrawFromSenderState:       types.DepositToBankState,
			types.DepositToBankState:            types.VerifyBankingTransactionState,
			types.VerifyBankingTransactionState: types.CompletedState,
			types.RefundToSenderState:           types.RejectedState,
		},
		failureMap: map[types.State]types.State{
			types.ValidatingState:               types.RejectedState,
			types.WithDrawFromSenderState:       types.RejectedState,
			types.DepositToBankState:            types.RefundToSenderState,
			types.VerifyBankingTransactionState: types.RefundToSenderState,
			types.RefundToSenderState:           types.RejectedState,
		},
	}
}

func (t *stateController) NextState(state types.State, status types.Status) types.State {
	if status == types.CompleteStatus {
		return state
	}
	if status == types.SuccessStatus {
		return t.successMap[state]
	}
	return t.failureMap[state]
}

