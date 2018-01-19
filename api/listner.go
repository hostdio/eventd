package api

import (
	"context"

	"github.com/hostdio/eventd/eventkit"
)

type Listener interface {
	Listen(context.Context, func(context.Context, eventkit.Event)) error
}
