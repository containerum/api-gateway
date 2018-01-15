package datastore

import (
	"database/sql"
	"fmt"

	"git.containerum.net/ch/api-gateway/pkg/model"

	"github.com/jmoiron/sqlx"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file" //file driver for migrations

	log "github.com/Sirupsen/logrus"
)

type data struct {
	*sqlx.DB
}

var (
	//MigrationsPath path to migrations folder
	MigrationsPath = "file://pkg/store/migrations"
)

//New create Store interface
func New(config model.DatabaseConfig) (interface{}, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Address,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	))
	if err != nil {
		log.WithError(err).Fatal(ErrUnableConnectPostgres)
		return nil, ErrUnableConnectPostgres
	}
	if config.Migrations {
		if err := runMigrationUP(db.DB); err != nil {
			log.WithError(err).Fatal(ErrUnableRunMigrations)
			return nil, ErrUnableRunMigrations
		}
	}
	return &data{db}, nil
}

func runMigrationUP(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.WithError(err).Fatal(ErrUnableCreatePostgresInstance)
		return ErrUnableCreatePostgresInstance
	}
	m, err := migrate.NewWithDatabaseInstance(MigrationsPath, "postgres", driver)
	if err != nil {
		log.WithError(err).Fatal(ErrUnableCreateMigrationInstance)
		return ErrUnableCreateMigrationInstance
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.WithError(err).Fatal(ErrUnableRunUpMigration)
		return ErrUnableRunUpMigration
	}
	version, dirty, err := m.Version()
	if err != nil {
		log.WithError(err).Fatal(ErrUnableGetMigrationVersion)
		return ErrUnableGetMigrationVersion
	}
	log.WithFields(log.Fields{
		"Version": version,
		"Dirty":   dirty,
	}).Info("Migration")
	return nil
}
