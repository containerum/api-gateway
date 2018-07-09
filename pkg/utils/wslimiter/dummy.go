package wslimiter

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type DummyLimiter struct {
	Upgrader *websocket.Upgrader
}

func (d *DummyLimiter) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	conn, err := d.Upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Conn:    conn,
		limiter: d,
	}, nil
}

func (d *DummyLimiter) untrack(*Conn) {}
