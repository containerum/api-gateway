package store

import (
	"context"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
)

type Store interface {
	Migrate(arg ...string) (string, error)

	//TestSelect create simple query in DB for check connection
	TestSelect() error

	//GetRouter gets router by unique uuid
	GetRouter(string) (*model.Router, error)

	//GetRoutesList gets all active routes
	GetRoutesList() (*[]model.Router, error)
}

func Migrate(c context.Context, arg ...string) (string, error) {
	return FromContext(c).Migrate(arg...)
}

//TestSelect create simple query in DB for check connection
func TestSelect(c context.Context) error {
	return FromContext(c).TestSelect()
}

//GetRouter gets router by unique uuid
func GetRouter(c context.Context, id string) (*model.Router, error) {
	return FromContext(c).GetRouter(id)
}

//GetRoutesList gets all active routes
func GetRoutesList(c context.Context) (*[]model.Router, error) {
	return FromContext(c).GetRoutesList()
}
