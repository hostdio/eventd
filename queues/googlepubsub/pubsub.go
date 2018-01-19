package googlepubsub

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hostdio/eventd/eventkit"

	"cloud.google.com/go/pubsub"
)

var (
	ErrMissingTopic = errors.New("Missing topic")
)

func New(ctx context.Context, projectID, topicID, subscriptionID string) (*Pubsub, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		client.Close()
		return nil, err
	} else if !exists {
		client.Close()
		return nil, ErrMissingTopic
	}
	subscription := client.Subscription(subscriptionID)
	return &Pubsub{
		client:       client,
		topic:        topic,
		subscription: subscription,
	}, nil
}

type Pubsub struct {
	client       *pubsub.Client
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

func (p Pubsub) Close() error {
	return p.client.Close()
}

func (p Pubsub) Listen(ctx context.Context, handler func(context.Context, eventkit.Event)) error {
	err := p.subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var event eventkit.Event
		if unmarshallErr := json.Unmarshal(m.Data, &event); unmarshallErr != nil {
			panic(unmarshallErr)
		}
		handler(ctx, event)
		m.Ack()
	})

	return err
}
