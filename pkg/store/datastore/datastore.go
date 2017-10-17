package datastore

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/pg"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store"
	migrate "bitbucket.org/exonch/ch-gateway/pkg/store/datastore/migrations"
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

func (db *datastore) Migrate(arg ...string) (string, error) {
	var answer string
	oldVersion, newVersion, err := migrate.RunMigration(db.DB, arg...)
	if err != nil {
		log.WithField("err", err.Error()).Fatal("Migration failed")
		return "", err
	}
	if newVersion != oldVersion {
		answer = fmt.Sprintf("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		answer = fmt.Sprintf("version is %d", oldVersion)
	}

	log.WithFields(log.Fields{
		"OldVersion": oldVersion,
		"NewVersion": newVersion,
		"Args":       arg,
		"Answer":     answer,
	}).Debug("Migration")

	return answer, nil
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
	log.Debug("New PG Connection")
	return nil
}
