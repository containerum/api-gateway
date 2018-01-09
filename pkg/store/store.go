package store

import (
	"context"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/store/datastore"

	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres driver
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
	FindListener(l *model.Listener) (*model.Listener, error)
	GetListenerList(l *model.Listener) (*[]model.Listener, error)
	UpdateListener(l *model.Listener, utype int) error
	CreateListener(l *model.Listener) (*model.Listener, error)
	DeleteListener(id string) error
	/* Group */
	GetGroupList(g *model.Group) (*model.Group, error)
	CreateGroup(g *model.Group) (*model.Group, error)
}

//New create new Store interface for working with data
func New(config model.DatabaseConfig) (Store, error) {
	//IDEA Connection pool
	db, err := newConnection(config)
	if config.Debug {
		db.Debug()
		db.LogMode(true)
	}
	if config.SafeMigration {
		db.AutoMigrate(&model.Role{}, &model.Group{}, &model.Plugin{}, &model.Listener{}) //"Safety" migrations
	}
	return datastore.New(db).(Store), err
}

// FromContext returns the Store associated with this context.
func FromContext(c context.Context) Store {
	return c.Value(key).(Store)
}

func GetListener(c context.Context, id string) (*model.Listener, error) {
	return FromContext(c).GetListener(id)
}

func FindListener(c context.Context, l *model.Listener) (*model.Listener, error) {
	return FromContext(c).FindListener(l)
}

func GetListenerList(c context.Context, l *model.Listener) (*[]model.Listener, error) {
	return FromContext(c).GetListenerList(l)
}

func UpdateListener(c context.Context, l *model.Listener, utype int) error {
	return FromContext(c).UpdateListener(l, utype)
}

func CreateListener(c context.Context, l *model.Listener) (*model.Listener, error) {
	return FromContext(c).CreateListener(l)
}

func DeleteListener(c context.Context, id string) error {
	return FromContext(c).DeleteListener(id)
}

func GetGroupList(c context.Context, g *model.Group) (*model.Group, error) {
	return FromContext(c).GetGroupList(g)
}

func CreateGroup(c context.Context, g *model.Group) (*model.Group, error) {
	return FromContext(c).CreateGroup(g)
}

func Init(c context.Context) error {
	return FromContext(c).Init()
}

func Version(c context.Context) (int, error) {
	return FromContext(c).Version()
}

func Up(c context.Context) (int, error) {
	return FromContext(c).Up()
}

func Down(c context.Context) (int, error) {
	return FromContext(c).Down()
}
