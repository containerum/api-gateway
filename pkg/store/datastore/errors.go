package datastore

import "errors"

var (
	ErrUnableConnectPostgres = errors.New("Unable to connect postgresql database")

	ErrUnableRunMigrations           = errors.New("Unable to run migrations")
	ErrUnableCreatePostgresInstance  = errors.New("Unable to create postgres instance for migrations")
	ErrUnableCreateMigrationInstance = errors.New("Unable to create migration instance")
	ErrUnableRunUpMigration          = errors.New("Unable to run UP migration")
	ErrUnableGetMigrationVersion     = errors.New("Unable to get migration version")

	ErrUnableGetListener    = errors.New("Unable to get listener")
	ErrUnableGetListeners   = errors.New("Unable to get listeners")
	ErrUnableCreateListener = errors.New("Unable to create listener")
	ErrUnableGetListenerID  = errors.New("Unable to get listener id")
	ErrUnableDeleteListener = errors.New("Unable to delete listener")
	ErrUnableUpdateListener = errors.New("Unable to update listener")

	ErrUnableGetGroups   = errors.New("Unable to get groups")
	ErrUnableCreateGroup = errors.New("Unable to create group")

	ErrNoRows = errors.New("No rows")

	ErrUnableScanListener   = errors.New("Unable to scan Listener struct")
	ErrUnableScanListenerID = errors.New("Unable to scan Listener id")

	ErrUnableScanGroup   = errors.New("Unable to scan Group struct")
	ErrUnableScanGroupID = errors.New("Unable to scan Group id")
)
