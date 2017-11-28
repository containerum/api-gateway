package main

import (
	"os"
	"time"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"bitbucket.org/exonch/ch-gateway/pkg/store/datastore"
	"bitbucket.org/exonch/ch-gateway/pkg/storex"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/urfave/cli"

	"bitbucket.org/exonch/ch-gateway/pkg/router"
	log "github.com/Sirupsen/logrus"
)

//setup DB and migration imp. in Store
func setupStore(c *cli.Context) store.Store {

	_, err := storex.New(model.DatabaseConfig{
		User:     c.String("pg-user"),
		Password: c.String("pg-password"),
		Database: c.String("pg-database"),
		Address:  c.String("pg-address"),
		Port:     c.String("pg-port"),
	})

	st := datastore.New(model.DatabaseConfig{
		User:     c.String("pg-user"),
		Password: c.String("pg-password"),
		Database: c.String("pg-database"),
		Address:  c.String("pg-address") + ":" + c.String("pg-port"),
	})

	if err != nil {
		panic(err)
	}

	if c.Bool("migrate") {
		answer, err := st.Migrate(c.Args()...)
		if err != nil {
			log.WithField("Error", err.Error()).Error("Migration failed")
		} else {
			log.WithField("Answer", answer).Info("Migration ok")
		}
		os.Exit(2)
	}

	return st
}

func setupRouters(r *router.Router, s store.Store) {
	routeList, err := s.GetRoutesList()
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error("GetRouteList failed id setupRouters")
	}
	for _, m := range *routeList {
		r.AddRoute(&m)
	}
}

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
