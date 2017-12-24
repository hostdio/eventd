package httpserver

import (
	"github.com/hostdio/eventd/api"
	"cloud.google.com/go/pubsub"

	"context"
	"encoding/json"
)

// TODO add interface in place of pubsub.Subscription
func startPersister(ctx context.Context, sub *pubsub.Subscription, persister api.Persister) error {

	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var event api.PublishedEvent
		if unmarshallErr := json.Unmarshal(m.Data, &event); unmarshallErr != nil {
			panic(unmarshallErr)
		}
		if persistErr := persister.Store(ctx, event); persistErr != nil {
			panic(persistErr)
		}
	})

	if err != nil {
		panic(err)
	}
	return nil
}
