package main

import (
	"context"
)

type (
	Handler interface {
		Handle(ctx context.Context, transID int64) error
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
			c.handler.Handle(ctx, transID)
		}
	}()
}