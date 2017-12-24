package api

import (
	"context"
	"encoding/json"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type BaseEvent struct {
	ID        string `validate:"required,lt=2000" json:"id"`
	Type      string `validate:"required,lt=2000" json:"type"`
	Version   string `validate:"required,lt=2000" json:"version"`
	Timestamp time.Time  `validate:"required" json:"timestamp"`
	Payload   string `json:"payload"`
	Source    string `validate:"required,lt=2000" json:"source"`
}

// PublishEvent represents the Event used to publish
type PublishEvent struct {
	*BaseEvent
}

// JSON turns PublishEvent into a JSON object
func (e PublishEvent) JSON() []byte {
	byt, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return byt
}

func (e PublishEvent) Validate() error {
	return validate.Struct(e)
}

func (e PublishEvent) Received() PublishedEvent {
	return PublishedEvent{
		PublishEvent: &e,
		ReceivedTimestamp: time.Now().UTC(),
	}
}

// PublishedEvent represents an event that has been received and pushed
type PublishedEvent struct {
	*PublishEvent
	ReceivedTimestamp time.Time `validate:"required"`
}

// JSON turns PublishedEvent into a JSON object
func (e PublishedEvent) JSON() []byte {
	byt, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return byt
}



// Publisher is the interface for implementing evnent publishers
type Publisher interface {
	Publish(ctx context.Context, event PublishEvent) (string, error)
}
