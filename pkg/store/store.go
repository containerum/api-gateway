package store

import (
	"context"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	uuid "github.com/satori/go.uuid"
)

type Store interface {
	//TestSelect create simple query in DB for check connection
	TestSelect() error

	//GetRouter gets router by unique uuid
	GetRouter(uuid.UUID) (*model.Router, error)

	//GetRoutesList gets all active routes
	GetRoutesList() ([]*model.Router, error)
}

//TestSelect create simple query in DB for check connection
func TestSelect(c context.Context) error {
	return FromContext(c).TestSelect()
}

//GetRouter gets router by unique uuid
func GetRouter(c context.Context, id uuid.UUID) (*model.Router, error) {
	return FromContext(c).GetRouter(id)
}

//GetRoutesList gets all active routes
func GetRoutesList(c context.Context) ([]*model.Router, error) {
	return FromContext(c).GetRoutesList()
}
