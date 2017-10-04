package store

import (
	"context"
)

const key = "store"

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Store {
	return c.Value(key).(Store)
}
