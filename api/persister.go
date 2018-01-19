package api

import (
	"context"
	"time"

	"github.com/hostdio/eventd/eventkit"
)

// PersistedEvent represent the event that has been persisted
type PersistedEvent struct {
	*eventkit.Event
	StoredTimestamp time.Time `validate:"required"`
}

type Persister interface {
	Store(context.Context, eventkit.Event) error
}
