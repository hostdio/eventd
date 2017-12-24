package httpserver

import (
	"github.com/hostdio/eventd/api"
	"context"
)

func startPersister(ctx context.Context, listener api.Listener, persister api.Persister) error {

	err := listener.Listen(ctx, func(ctx context.Context, event api.PublishedEvent) {
		if persistErr := persister.Store(ctx, event); persistErr != nil {
			panic(persistErr)
		}
	})

	return err
}
