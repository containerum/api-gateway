package main

import (
	"fmt"
	"os"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	"bitbucket.org/exonch/ch-gateway/pkg/store/datastore"
	"github.com/urfave/cli"
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
			fmt.Print(err.Error())
		}
	}

	if c.Bool("migrate") {
		answer, err := st.Migrate(c.Args()...)
		if err != nil {
			fmt.Print(err.Error())
		} else {
			fmt.Print(answer)
		}
		os.Exit(2)
	}

	return st
}
