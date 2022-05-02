package main

import (
	"context"
	"ewallet/modules/transaction"
)

type (
	transactionChannel struct {
		transChannel chan int64
		consumer     transaction.StateMachine
	}
)

func NewPublisher(transChannel chan int64) *transactionChannel {
	return &transactionChannel{
		transChannel: transChannel,
	}
}

func (c *transactionChannel) Publish(ctx context.Context, transId int64) error {
	c.transChannel <- transId
	return nil
}