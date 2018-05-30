package wslimiter

import (
	"errors"
	"net/http"
)

var ErrLimitReached = errors.New("connection limit reached")

type LimitedUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error)

	untrack(*Conn)
}
