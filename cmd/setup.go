package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server"
	"git.containerum.net/ch/api-gateway/pkg/store"

	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"
	"git.containerum.net/ch/grpc-proto-files/auth"
	rate "git.containerum.net/ch/ratelimiter"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
)

func setupServer(c *cli.Context) *server.Server {
	serv := server.New(time.Second * 5)
	if conf.Store.Enable {
		serv.RegisterStore(setupStore(c))
	}
	if conf.Auth.Enable {
		serv.RegisterAuth(setupAuth(c))
	}
	if conf.Rate.Enable {
		serv.RegisterRatelimiter(setupRatelimiter(c))
	}
	if conf.Clickhouse.Enable {
		serv.RegisterClickhouseLogger(setupClickhouseLogger(c))
	}
	return serv
}

func setupLogger(c *cli.Context) {
	if c.Bool("debug") {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
		log.Debug("Application running in Debug mode")
		return
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

//setup DB and migration imp. in Store
func setupStore(c *cli.Context) *store.Store {
	log.WithFields(log.Fields{
		"PG_USER":       c.String("pg-user"),
		"PG_PASSWORD":   c.String("pg-password"),
		"PG_DATABASE":   c.String("pg-database"),
		"PG_ADDRESS":    c.String("pg-address"),
		"PG_PORT":       c.String("pg-port"),
		"PG_MIGRATIONS": c.Bool("pg-migrations"),
		"PG_DEBUG":      c.Bool("pg-debug"),
	}).Debug("Setup DB connection")

	st, err := store.New(model.DatabaseConfig{
		User:       c.String("pg-user"),
		Password:   c.String("pg-password"),
		Database:   c.String("pg-database"),
		Address:    c.String("pg-address"),
		Port:       c.String("pg-port"),
		Debug:      c.Bool("pg-debug"),
		Migrations: c.Bool("pg-migrations"),
	})
	if err != nil {
		panic(err)
	}
	return &st
}

func setupAuth(c *cli.Context) *auth.AuthClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	addr := fmt.Sprintf("%s:%s", c.String("grpc-auth-address"), c.String("grpc-auth-port"))

	con, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("Dial to auth grpc failed")
	}

	client := auth.NewAuthClient(con)
	return &client
}

func setupRatelimiter(c *cli.Context) *rate.PerIPLimiter {
	maxRate, err := strconv.Atoi(c.String("rate-limit"))
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("Rate limit parse error")
	}

	ratelimiter, err := rate.NewPerIPLimiter(&rate.PerIPLimiterConfig{
		RedisAddress:  c.String("redis-address"),
		RedisPassword: c.String("redis-password"),
		RedisDB:       1,
		Timeout:       time.Second,
		RateLimit:     maxRate,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("Connection redis limiter failed")
	}
	return ratelimiter
}

func setupStatsd(c *cli.Context) *statsd.Statter {
	std, err := statsd.NewBufferedClient(
		c.String("statsd-address"),
		c.String("statsd-prefix"),
		time.Microsecond*time.Duration(c.Int("statsd-buffer-time")),
		0,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"Err":         err,
			"Address":     c.String("statsd-address"),
			"Prefix":      c.String("statsd-prefix"),
			"Buffer-Time": c.Int("statsd-buffer-time"),
		}).Warning("Setup Statsd failed")
	}
	return &std
}

func setupClickhouseLogger(c *cli.Context) *clickhouse.LogClient {
	client, err := clickhouse.OpenConenction(c.String("clickhouse-logger"))
	if err != nil {
		log.WithFields(log.Fields{
			"Err":     err,
			"Address": c.String("clickhouse-logger"),
		}).Warning("Setup Clickhouse Logger failed")
	}
	return client
}

func setupTSL(c *cli.Context) (certFile string, keyFile string, err error) {
	if _, e := os.Stat(c.String("tls-cert")); os.IsNotExist(e) {
		log.WithError(err).Error("Unable to open cert.pem")
		err = errors.New("Unable to open cert.pem")
		return
	}
	if _, e := os.Stat(c.String("tls-key")); os.IsNotExist(e) {
		log.WithError(err).Error("Unable to open key.pem")
		err = errors.New("Unable to open key.pem")
		return
	}
	return c.String("tls-cert"), c.String("tls-key"), nil
}
