package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/store"

	clickhouse "git.containerum.net/ch/api-gateway/pkg/utils/clickhouselog"
	"git.containerum.net/ch/grpc-proto-files/auth"
	rate "git.containerum.net/ch/ratelimiter"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
)

//setup DB and migration imp. in Store
func setupStore(c *cli.Context) *store.Store {
	st, err := store.New(model.DatabaseConfig{
		User:          c.String("pg-user"),
		Password:      c.String("pg-password"),
		Database:      c.String("pg-database"),
		Address:       c.String("pg-address"),
		Port:          c.String("pg-port"),
		Debug:         c.Bool("pg-debug"),
		SafeMigration: c.Bool("pg-safe-migration"),
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

// _, err = client.CheckToken(context.Background(), &auth.CheckTokenRequest{
// 	AccessToken: "sss",
// 	UserAgent:   "chrome",
// 	FingerPrint: "x1x2",
// 	UserIp:      "192.168.0.1",
// })
// if err != nil {
// 	log.WithFields(log.Fields{
// 		"Err": grpc.ErrorDesc(err),
// 	}).Error("CheckToken error")
// }

// res, err := client.CreateToken(context.Background(), &auth.CreateTokenRequest{
// 	UserAgent:   "Mozilla/5.0 (X11; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0",
// 	Fingerprint: "550e8400-e29b-41d4-a716-446655440000",
// 	UserId:      &common.UUID{Value: "550e8400-e29b-41d4-a716-446655440000"},
// 	UserIp:      "127.0.0.1",
// 	UserRole:    auth.Role_USER,
// })
// if err != nil {
// 	log.WithFields(log.Fields{
// 		"Err": grpc.ErrorDesc(err),
// 	}).Error("CreateToken error")
// }
