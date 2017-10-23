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

	//GetRoutesListActivation return list of all routers
	GetRoutesListByActivation(active bool) (*[]model.Router, error)

	//AddRouter create new router
	AddRouter(r *model.Router) error
}

//Migrate run migrations
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

//GetRoutesListByActivation return list of all routers
func GetRoutesListByActivation(c context.Context, active bool) (*[]model.Router, error) {
	return FromContext(c).GetRoutesListByActivation(active)
}

//AddRouter create new router
func AddRouter(c context.Context, r *model.Router) error {
	return FromContext(c).AddRouter(r)
}
