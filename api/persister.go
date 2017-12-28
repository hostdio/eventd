package api

import (
	"time"
	"context"
)

// PersistedEvent represent the event that has been persisted
type PersistedEvent struct {
	*PublishedEvent
	StoredTimestamp time.Time `validate:"required"`
}

type Persister interface {
	Store(context.Context, PublishedEvent) error
}



