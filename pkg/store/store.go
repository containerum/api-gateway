package store

import (
	"context"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/store/datastore"

	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres driver
)

const key = "store"

//Store impl functions for working with data
type Store interface {
	/*Migration */
	Init() error
	Version() (int, error)
	Up() (int, error)
	Down() (int, error)
	/* Listener */
	GetListener(id string) (*model.Listener, error)
	FindListener(l *model.Listener) (*model.Listener, error)
	GetListenerList(l *model.Listener) (*[]model.Listener, error)
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
