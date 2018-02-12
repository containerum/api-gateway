package datastore

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	migrationsPath = "file://pkg/store/migrations"
)

//ListenerGroup struct for scan listener and group in one scan
type ListenerGroup struct {
	model.Listener
	GroupName      string    `db:"group_name"`
	GroupCreatedAt time.Time `db:"group_created_at"`
	GroupUpdatedAt time.Time `db:"group_updated_at"`
	GroupActive    bool      `db:"group_active"`
}

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
	log.WithField("Path", migrationsPath).Debug("Migrations")
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
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

func initRows() *sqlx.Rows {
	return &sqlx.Rows{}
}

func scanListenerGroup(rows *sqlx.Rows) (*model.Listener, error) {
	var localListener ListenerGroup
	if err := rows.StructScan(&localListener); err != nil {
		log.WithError(err).Warn(ErrUnableScanListener)
		return nil, ErrUnableScanListener
	}
	listener := &model.Listener{
		DefaultModel: model.DefaultModel{
			ID:        localListener.ID,
			CreatedAt: localListener.CreatedAt,
			UpdatedAt: localListener.UpdatedAt,
		},
		Name:       localListener.Name,
		OAuth:      localListener.OAuth,
		Active:     localListener.Active,
		GroupRefer: localListener.GroupRefer,
		Group: model.Group{
			DefaultModel: model.DefaultModel{
				ID:        localListener.GroupRefer,
				CreatedAt: localListener.GroupCreatedAt,
				UpdatedAt: localListener.GroupUpdatedAt,
			},
			Name: localListener.GroupName,
		},
		StripPath:   localListener.StripPath,
		ListenPath:  localListener.ListenPath,
		UpstreamURL: localListener.UpstreamURL,
		Method:      localListener.Method,
		Roles:       localListener.Roles,
	}
	listener.SetRoles(parseRoles(localListener.Roles))
	return listener, nil
}

func parseRoles(br []byte) []string {
	var roles []string
	str := string(br)
	str = strings.TrimLeft(str, "{")
	str = strings.TrimRight(str, "}")
	roles = strings.Split(str, ",")
	return roles
}
