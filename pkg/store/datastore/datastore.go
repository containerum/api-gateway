package datastore

import (
	"fmt"

	"github.com/go-pg/pg"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
)

// datastore is an implementation of a model.Store
type datastore struct {
	*pg.DB
	config *model.DatabaseConfig
}

// New creates a database connection for the given driver and datasource
// and returns a new Store.
func New(config model.DatabaseConfig) store.Store {
	return &datastore{
		DB:     open(&config),
		config: &config,
	}
}

//TestSelect create simple query in DB for check connection
func (db *datastore) TestSelect() error {
	_, err := db.Exec("SELECT 1")
	return err
}

func open(config *model.DatabaseConfig) *pg.DB {
	return pg.Connect(&pg.Options{
		User:      config.User,
		Password:  config.Password,
		Database:  config.Database,
		Addr:      config.Address,
		OnConnect: onConnect,
	})
}

func onConnect(con *pg.DB) error {
	fmt.Print("Connected")
	return nil
}
