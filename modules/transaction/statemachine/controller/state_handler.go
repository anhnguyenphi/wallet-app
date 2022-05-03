package controller

import (
	"context"
	"ewallet/modules/share/dbclient"
	"ewallet/modules/transaction/statemachine/types"
)

const (
	selectForUpdateTransactionSql     = `SELECT id, type, from_wallet_id, from_card_id, to_wallet_id, to_card_id, external_tx_id, amount, currency, state, created_at, updated_at from transactions where id = ? FOR UPDATE`
	updateStateOfTransactionForUpdate = `UPDATE transactions SET state = ?, external_tx_id = ? WHERE id = ?`
)

type (

	MachineHandler interface {
		Handle(ctx context.Context, transID int64) (types.Status, error)
	}

	stateMachineHandler struct {
		transactionDB      dbclient.DB
		stateHandlerGetter types.StateHandlerGetter
		stateController    types.StateController
	}
)

func NewStateMachineHandler(transactionDB dbclient.DB, handlerGetter types.StateHandlerGetter, stateController types.StateController) MachineHandler {
	return &stateMachineHandler{
		transactionDB:      transactionDB,
		stateHandlerGetter: handlerGetter,
		stateController: stateController,
	}
}

func (r *stateMachineHandler) Handle(ctx context.Context, transID int64) (types.Status, error) {
	tx, err := r.transactionDB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var trans types.TransDAO
	result := tx.QueryRowContext(ctx, selectForUpdateTransactionSql, transID)

	if err := trans.Scan(result); err != nil {
		return "", err
	}

	stateHandler, err := r.stateHandlerGetter.GetHandler(trans.State)
	if err != nil {
		return "", err
	}
	stateTx, err := stateHandler.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer stateTx.Rollback()

	status, err := stateHandler.Handle(ctx, stateTx, &trans)
	if err != nil {
		return "", err
	}
	nextState := r.stateController.NextState(trans.State, status)
	if _, err = tx.ExecContext(ctx, updateStateOfTransactionForUpdate,
		nextState, trans.ExTransactionID.String, transID); err != nil {
		return "", err
	}

	if err := stateTx.Commit(); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return status, nil
}
