package googlepubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/hostdio/eventd/api"
	"errors"
)

var (
	ErrMissingTopic = errors.New("Missing topic")
)

func New(ctx context.Context, projectID, topicID string) (*Pubsub, error) {
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
	return &Pubsub{
		client:client,
		topic: topic,
	}, nil
}

type Pubsub struct {
	client *pubsub.Client
	topic *pubsub.Topic
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