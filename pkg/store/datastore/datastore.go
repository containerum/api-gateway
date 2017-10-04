package datastore

import (
	"database/sql"

	"bitbucket.org/exonch/ch-gateway/pkg/store"
)

// datastore is an implementation of a model.Store built on top
// of the sql/database driver with a relational database backend.
type datastore struct {
	*sql.DB

	driver string
	config string
}

// New creates a database connection for the given driver and datasource
// and returns a new Store.
func New(driver, config string) store.Store {
	// return &datastore{
	// DB:     open(driver, config),
	// driver: driver,
	// config: config,
	// }

	return &datastore{}
}
