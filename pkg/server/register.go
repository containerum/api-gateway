package server

import (
	"git.containerum.net/ch/api-gateway/pkg/store"
	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/ratelimiter"

	log "github.com/Sirupsen/logrus"
)

//RegisterStore registre store interface in router
func (s *Server) RegisterStore(st *store.Store) {
	log.WithField("Store", st).Debug("Register store in server")
	s.storeClient = st
}

//RegisterAuth registre auth interface in router
func (s *Server) RegisterAuth(c *auth.AuthClient) {
	log.WithField("Auth", c).Debug("Register auth in server")
	s.authClient = c
}

//RegisterRatelimiter registre statsd interface in router
func (s *Server) RegisterRatelimiter(l *ratelimiter.PerIPLimiter) {
	log.WithField("Rate", l).Debug("Register rate limiter in server")
	s.rateClient = l
}

//RegisterClickhouseLogger registre clickhouse logger client
func (s *Server) RegisterClickhouseLogger(c *clickhouse.LogClient) {
	log.WithField("Logger", c).Debug("Register clickhouse logger in server")
	s.clickhouseLoggerClient = c
}
