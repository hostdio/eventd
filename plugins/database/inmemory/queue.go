package inmemory

import (
	"context"

	"github.com/hostdio/eventd/eventkit"
)

func NewQueue() Queue {
	return Queue{
		c: make(chan eventkit.Event),
	}
}

type Queue struct {
	c chan eventkit.Event
}

func (p Queue) Listen(ctx context.Context, handler func(context.Context, eventkit.Event)) error {
	select {
	case event := <-p.c:
		handler(ctx, event)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
