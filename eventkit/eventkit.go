package eventkit

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// Event is the common event structure expected to be used
// by event producers and consumers
type Event struct {
	Namespace string    `json:"namespace"`
	Type      string    `json:"type"`
	ID        string    `json:"id"`
	Version   string    `json:"version"`
	Source    string    `json:"source"`
	Produced  time.Time `json:"produced"`
	Data      JSONData  `json:"data"`
	Metadata  JSONData  `json:"metadata"`
}

// JSON serializes itself into a JSON object
func (e Event) JSON() []byte {
	byt, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return byt
}

// JSONData is the abstraction of a JSON data structure
type JSONData map[string]interface{}

// JSON serializes itself into a JSON object
func (d JSONData) JSON() []byte {
	byt, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return byt
}

// Consumer is the interface used to implement an event consumer
type Consumer interface {
	Consume(ctx context.Context, consume func(Event) error) error
	Close() error
}

// Producer is the interface used to implement an event producer
type Producer interface {
	Publish(context.Context, Event) error
	Close() error
}

// PubsubProducer is an example of a producer
type PubsubProducer struct {
	ps      *pubsub.Client
	topicID string
}

// Publish publishes an event to pubsub
func (p PubsubProducer) Publish(ctx context.Context, e Event) error {
	topic := p.ps.Topic(p.topicID)
	msg := pubsub.Message{
		Data: e.JSON(),
	}
	res := topic.Publish(ctx, &msg)
	_, err := res.Get(ctx)
	return err
}

// Close closes the Producer connection to Pubsub
func (p PubsubProducer) Close() error {
	return p.ps.Close()
}

// PubsubConsumer consumes the pubsub events made available to a subscription
type PubsubConsumer struct {
	ps              *pubsub.Client
	subscriptionID  string
	receiveSettings pubsub.ReceiveSettings
}

// Consume takes incoming events and consumes them asynchronously
// the consume function returns an Event to the consumer. If the function
// returns an error the system will not acknowledge the Event and it will
// be put back on the event queue. If it returns nil Consume will
// acknowledge the messag for you
func (c PubsubConsumer) Consume(ctx context.Context, consume func(Event) error) error {
	sub := c.ps.Subscription(c.subscriptionID)
	sub.ReceiveSettings = c.receiveSettings
	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		byt := msg.Data
		var event Event
		err := json.Unmarshal(byt, &event)
		if err != nil {
			// assume that the event is broken
			log.Println("Received a broken event", string(byt))
			msg.Ack()
			return
		}

		err = consume(event)
		if err != nil {
			log.Println("Errors are not handled as of now: ", err.Error())
			msg.Nack()
		}
		msg.Ack()
	})
	return err
}

// Close closes the Producer connection to Pubsub
func (c PubsubConsumer) Close() error {
	return c.ps.Close()
}
