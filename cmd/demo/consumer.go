package main

import (
	"context"
	"ewallet/modules/transaction/statemachine/controller"
	"ewallet/modules/transaction/statemachine/types"
)

type (
	Handler interface {
		Handle(ctx context.Context, transID int64) (types.Status, error)
	}
	consumer struct {
		transChannel chan int64
		handler Handler
	}
)

func NewConsumer(transChannel chan int64, handler Handler) *consumer  {
	return &consumer{
		transChannel: transChannel,
		handler: handler,
	}
}

func (c *consumer) Run(ctx context.Context) {
	go func() {
		for transID := range c.transChannel {
			status, err := c.handler.Handle(ctx, transID)
			if err != nil || status != controller.CompleteStatus {
				c.transChannel <- transID
			}
		}
	}()
}