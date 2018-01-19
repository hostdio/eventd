package api

import (
	"context"
	"time"

	"github.com/hostdio/eventd/eventkit"
)

type Scanner interface {
	Scan(context.Context, time.Time, int) ([]eventkit.Event, error)
}
