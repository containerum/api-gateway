package store

import (
	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/store/datastore"
)

const key = "store"

//Store impl functions for working with data
type Store interface {
	/*Migrations */
	Init() error
	Version() (int, error)
	Up() (int, error)
	Down() (int, error)
	/* Listener */
	GetListener(id string) (*model.Listener, error)
	GetListenerList(l *model.Listener) (*[]model.Listener, error)
	UpdateListener(l *model.Listener, utype model.ListenerUpdateType) error
	CreateListener(l *model.Listener) (*model.Listener, error)
	DeleteListener(id string) error
	/* Group */
	GetGroupList(g *model.Group) (*model.Group, error)
	CreateGroup(g *model.Group) (*model.Group, error)
}

//New create new Store interface for working with data
func New(config model.DatabaseConfig) (Store, error) {
	db, err := datastore.New(config)
	return db.(Store), err
}
