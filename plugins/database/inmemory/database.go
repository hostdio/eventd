package inmemory

import (
	"sync"
	"context"
	"github.com/hostdio/eventd/api"
	"time"
)

func NewDatabase() Database {
	return Database{
		lock: new(sync.RWMutex),
		data: make(map[string]api.PersistedEvent),
	}
}

type Database struct {
	lock *sync.RWMutex
	data map[string]api.PersistedEvent
}

func (c Database) Store(ctx context.Context, event api.PublishedEvent) error {
	c.lock.Lock()
	persistedEvent := api.PersistedEvent{
		PublishedEvent: &event,
		StoredTimestamp: time.Now(),
	}
	c.data[event.ID] = persistedEvent
	defer c.lock.Unlock()
	return nil
}

func (c Database) Scan(ctx context.Context, from time.Time, limit int) ([]api.PersistedEvent, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	events := []api.PersistedEvent{}
	for _, event := range c.data {
		if event.StoredTimestamp.Before(from) {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}