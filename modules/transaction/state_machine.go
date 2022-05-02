package transaction

import (
	"context"
	"ewallet/modules/share"
)

type (
	State string

	StateMachine interface {
		Handle(ctx context.Context, transID int64) error
	}

	stateHandler interface {
		Handle(ctx context.Context,tx share.TX, trans TransDAO) (State, error)
		Begin(ctx context.Context) (share.TX, error)
	}

	stateMachine struct {
		transactionDB      share.DB
		assetDB            share.DB
		stateHandlerGetter StateHandlerGetter
	}
)

func NewStateMachine(transactionDB ,assetDB share.DB, handlerGetter StateHandlerGetter) StateMachine {
	return &stateMachine{
		transactionDB:      transactionDB,
		assetDB:            assetDB,
		stateHandlerGetter: handlerGetter,
	}
}

func (r *stateMachine) Handle(ctx context.Context, transID int64) error {
	tx, err := r.transactionDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var trans TransDAO
	result := tx.QueryRowContext(ctx, selectForUpdateTransactionSql, transID)

	if err := result.Scan(&trans); err != nil {
		return err
	}

	stateHandler, err := r.stateHandlerGetter.GetHandler(trans.State)
	if err != nil {
		return err
	}
	stateTx, err := stateHandler.Begin(ctx)
	if err != nil {
		return err
	}
	defer stateTx.Rollback()

	nextState, err := stateHandler.Handle(ctx, stateTx, trans)
	if nextState == trans.State && err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, updateStateOfTransactionForUpdate, nextState); err != nil {
		return err
	}

	if err := stateTx.Commit(); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}