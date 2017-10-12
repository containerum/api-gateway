package datastore

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/pg"
	uuid "github.com/satori/go.uuid"

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
	st := &datastore{
		DB:     open(&config),
		config: &config,
	}

	rg := model.Group{
		ID: "550e8400-e29b-41d4-a716-446655440000",
	}

	err := st.Model(&rg).Select()
	if err != nil {
		fmt.Printf("Select err: %v", err)
	}

	r := &model.Router{
		ID:      uuid.NewV4().String(),
		Group:   &rg,
		OAuth:   true,
		Active:  true,
		Created: time.Now(),
	}

	_, err = st.Model(r).Insert()
	if err != nil {
		fmt.Printf("Insert err: %v", err)
	}

	return st
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
