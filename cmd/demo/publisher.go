package main

import (
	"context"
	"ewallet/modules/transaction/statemachine/controller"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	transactionChannel struct {
		transferChannel chan int64
		depositChannel  chan int64
		withdrawChannel chan int64
		consumer        controller.MachineHandler
	}
)

func NewPublisher(transChannel chan int64, depositChannel chan int64, withdrawChannel chan int64) *transactionChannel {
	return &transactionChannel{
		transferChannel: transChannel,
		depositChannel: depositChannel,
		withdrawChannel: withdrawChannel,
	}
}

func (c *transactionChannel) Publish(ctx context.Context, transType types.Type, transId int64) error {
	switch transType {
	case controller.DepositType: c.depositChannel <- transId
	case controller.TransferType: c.transferChannel <- transId
	case controller.WithdrawType: c.withdrawChannel <- transId
	}
	return nil
}