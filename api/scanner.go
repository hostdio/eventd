package api

import (
	"time"
	"context"
)

type Scanner interface {
	Scan(context.Context, time.Time) []PersistedEvent
}
