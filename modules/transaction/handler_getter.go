package transaction

import (
	"errors"
	"ewallet/modules/share"
)

type (

	StateHandlerGetter interface {
		GetHandler(state State) (stateHandler, error)
	}

	handlerGetterImpl struct {
		assetDB share.DB
		handlerMap map[State]stateHandler
	}
)

var InvalidState = errors.New("invalid state")

func NewStateHandlerGetter(assetDB share.DB) StateHandlerGetter {
	return &handlerGetterImpl{
		assetDB: assetDB,
		handlerMap: map[State]stateHandler{
			ValidatingState: &validatingStateHandler{
				assetDB: assetDB,
			},
			WithDrawFromSenderState: &withdrawStateHandler{
				assetDB: assetDB,
			},
			DepositToReceiverState: &depositStateHandler{
				assetDB: assetDB,
			},
			ReversingToSenderState: &reversingStateHandler{
				assetDB: assetDB,
			},
		},
	}
}

func (h *handlerGetterImpl) GetHandler(state State) (stateHandler, error) {
	handler, ok := h.handlerMap[state]
	if !ok {
		return nil, InvalidState
	}
	return handler, nil
}


