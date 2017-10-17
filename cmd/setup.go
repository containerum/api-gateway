package main

import (
	"fmt"
	"os"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"bitbucket.org/exonch/ch-gateway/pkg/store/datastore"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

//setup DB and migration imp. in Store
func setupStore(c *cli.Context) store.Store {
	st := datastore.New(model.DatabaseConfig{
		User:     c.String("pg-user"),
		Password: c.String("pg-password"),
		Database: c.String("pg-database"),
		Address:  c.String("pg-address"),
	})

	if c.Bool("debug") {
		err := st.TestSelect()
		if err != nil {
			log.WithField("Error", err.Error()).Error("Test select failed")
		}
	}

	if c.Bool("migrate") {
		answer, err := st.Migrate(c.Args()...)
		if err != nil {
			fmt.Print(err.Error())
			log.WithField("Error", err.Error()).Error("Migration failed")
		} else {
			log.WithField("Answer", answer).Info("Migration ok")
		}
		os.Exit(2)
	}

	return st
}
