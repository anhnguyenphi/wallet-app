package state

import (
	"errors"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	handlerGetterImpl struct {
		handlerMap map[types.State]types.StateHandler
	}
)

var InvalidState = errors.New("invalid state")

func NewStateHandlerGetter(assetDB dbclient.DB, deposit BankDeposit, withdraw BankWithdraw, verifier BankTransactionVerifier) types.StateHandlerGetter {
	return &handlerGetterImpl{
		handlerMap: map[types.State]types.StateHandler{
			types.ValidatingState: &validatingStateHandler{
				card: assetDB,
			},
			types.WithDrawFromSenderState: &withdrawFromWalletStateHandler{
				assetDB: assetDB,
			},
			types.DepositToReceiverState: &depositToWalletStateHandler{
				assetDB: assetDB,
			},
			types.RefundToSenderState: &refundToWalletStateHandler{
				assetDB: assetDB,
			},
			types.DepositToBankState: &depositToBankStateHandler{
				bank: deposit,
			},
			types.WithDrawFromBankState: &withdrawFromBankStateHandler{
				bank: withdraw,
			},
			types.VerifyBankingTransactionState: &verifyBankingTrxStateHandler{
				bank: verifier,
			},
			types.CompletedState: &completeStateHandler{
			},
			types.RejectedState: &rejectedStateHandler{
			},
		},
	}
}

func (h *handlerGetterImpl) GetHandler(state types.State) (types.StateHandler, error) {
	handler, ok := h.handlerMap[state]
	if !ok {
		return nil, InvalidState
	}
	return handler, nil
}


