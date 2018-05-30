package wslimiter

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
)

type IPLookupFunc func(r *http.Request) string

type PerIPLimitedUpgrader struct {
	cache          *cache.Cache
	maxConnections uint
	upgrader       *websocket.Upgrader
	ipLookup       IPLookupFunc

	mu sync.RWMutex // to atomize upgrade and untrack operations
}

func NewPerIPLimiter(maxConnections uint, upgrader *websocket.Upgrader, ipLookup IPLookupFunc) *PerIPLimitedUpgrader {
	if ipLookup == nil {
		ipLookup = func(r *http.Request) string {
			return r.RemoteAddr
		}
	}
	return &PerIPLimitedUpgrader{
		cache:          cache.New(time.Hour, time.Minute),
		maxConnections: maxConnections,
		upgrader:       upgrader,
		ipLookup:       ipLookup,
	}
}

func (p *PerIPLimitedUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	ip := p.ipLookup(r)

	p.mu.RLock()
	connections, ok := p.cache.Get(ip)
	if !ok {
		connections = uint(0)
	}
	p.mu.RUnlock()

	if connections.(uint) >= p.maxConnections {
		return nil, ErrLimitReached
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	conn, err := p.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	p.cache.IncrementUint(ip, 1)

	return &Conn{
		Conn:     conn,
		limiter:  p,
		cacheKey: ip,
	}, nil
}

func (p *PerIPLimitedUpgrader) untrack(conn *Conn) {
	key := conn.cacheKey

	p.mu.RLock()
	connections, ok := p.cache.Get(key)
	if !ok {
		connections = uint(0)
	}
	p.mu.RUnlock()

	if connections.(uint) == 0 {
		return
	}

	p.mu.Lock()
	p.cache.DecrementUint(key, 1)
	p.mu.Unlock()
}
