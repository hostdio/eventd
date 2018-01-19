package inmemory

import (
	"context"
	"github.com/hostdio/eventd/api"
	"time"
)

func NewQueue() Queue {
	return Queue{
		c: make(chan api.PublishedEvent),
	}
}

type Queue struct {
	c chan api.PublishedEvent
}

func (in Queue) Publish(ctx context.Context, event api.PublishEvent) (string, error) {
	select {
		case in.c<-event.Received():
			return time.Now().String(), nil
		case <-ctx.Done():
			return "", ctx.Err()
	}
}

func (p Queue) Listen(ctx context.Context, handler api.EventHandler) error {
	select {
	case event:=<-p.c:
		handler(ctx, event)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
