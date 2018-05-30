package wslimiter

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn
	limiter     LimitedUpgrader
	cacheKey    string
	onceUntrack sync.Once
}

func (c *Conn) Close() error {
	if c.limiter != nil {
		c.onceUntrack.Do(func() {
			c.limiter.untrack(c)
		})
	}
	return c.Conn.Close()
}
