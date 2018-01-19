package httpserver

import (
	"context"

	"github.com/hostdio/eventd/api"
	"github.com/hostdio/eventd/eventkit"
)

func startPersister(ctx context.Context, listener api.Listener, persister api.Persister) error {

	err := listener.Listen(ctx, func(ctx context.Context, event eventkit.Event) {
		if persistErr := persister.Store(ctx, event); persistErr != nil {
			panic(persistErr)
		}
	})

	return err
}
