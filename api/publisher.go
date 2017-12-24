package api

import (
	"encoding/json"
	"time"
	"context"
)

// PublishEvent represents the Event used to publish
type PublishEvent struct {
	ID        [2000]byte
	Type      [2000]byte
	Version   [2000]byte
	Timestamp time.Time
	Payload   []byte
	Source    [2000]byte
}

// JSON turns PublishEvent into a JSON object
func (e PublishEvent) JSON() []byte {
	byt, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return byt
}

// PersistedEvent represent the event that has been persisted
type PersistedEvent struct {
	ID                [2000]byte
	Type              [2000]byte
	Version           [2000]byte
	Timestamp         time.Time
	Payload           []byte
	Source            [2000]byte
	ReceivedTimestamp time.Time
}

// Publisher is the interface for implementing evnent publishers
type Publisher interface {
	Publish(ctx context.Context, event PublishEvent) (string, error)
}
