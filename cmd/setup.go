package main

import (
	"time"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

//setup DB and migration imp. in Store
func setupStore(c *cli.Context) store.Store {
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
	return st
}

// func setupRouters(r *router.Router, s store.Store) {
// 	routeList, err := s.GetRoutesList()
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"Err": err,
// 		}).Error("GetRouteList failed id setupRouters")
// 	}
// 	for _, m := range *routeList {
// 		r.AddRoute(&m)
// 	}
// }

//TODO: Make statsd client Own interface and move it to Store statsd package
func setupStatsd(c *cli.Context) statsd.Statter {
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

	return std
}
