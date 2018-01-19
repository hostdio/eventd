package httpserver

import (
	"context"

	"github.com/hostdio/eventd/api"
	"github.com/hostdio/eventd/eventkit"
)

func startPersister(ctx context.Context, consumer eventkit.Consumer, persister api.Persister) error {

	err := consumer.Consume(ctx, func(event eventkit.Event) error {
		if persistErr := persister.Store(ctx, event); persistErr != nil {
			return persistErr
		}
		return nil
	})

	return err
}
