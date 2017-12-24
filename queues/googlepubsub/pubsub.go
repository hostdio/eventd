package googlepubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/hostdio/eventd/api"
	"errors"
	"encoding/json"
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
		client:client,
		topic: topic,
		subscription: subscription,
	}, nil
}

type Pubsub struct {
	client *pubsub.Client
	topic *pubsub.Topic
	subscription *pubsub.Subscription
}

func (p Pubsub) Close() error {
	return p.client.Close()
}

func (p Pubsub) Publish(ctx context.Context, event api.PublishEvent) (string, error) {
	msg := &pubsub.Message{
		Data: event.JSON(),
	}
	res := p.topic.Publish(ctx, msg)
	serverID, err := res.Get(ctx)
	return serverID, err
}

func (p Pubsub) Listen(ctx context.Context, handler api.EventHandler) error {
	err := p.subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var event api.PublishedEvent
		if unmarshallErr := json.Unmarshal(m.Data, &event); unmarshallErr != nil {
			panic(unmarshallErr)
		}
		handler(ctx, event)
		m.Ack()
	})

	return err
}