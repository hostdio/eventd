package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/hostdio/eventd/eventkit"
)

func NewDatabase() Database {
	return Database{
		lock: new(sync.RWMutex),
		data: make(map[string]eventkit.Event),
	}
}

type Database struct {
	lock *sync.RWMutex
	data map[string]eventkit.Event
}

func (c Database) Store(ctx context.Context, event eventkit.Event) error {
	c.lock.Lock()

	c.data[event.ID] = event
	defer c.lock.Unlock()
	return nil
}

func (c Database) Scan(ctx context.Context, from time.Time, limit int) ([]eventkit.Event, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	events := []eventkit.Event{}
	for _, event := range c.data {
		if event.Produced.Before(from) {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}
