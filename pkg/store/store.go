package store

import (
	"context"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	uuid "github.com/satori/go.uuid"
)

type Store interface {
	//GetRouter gets router by unique uuid
	GetRouter(uuid.UUID) (*model.Router, error)

	//GetRoutesList gets all active routes
	GetRoutesList() ([]*model.Router, error)
}

//GetRouter gets router by unique uuid
func GetRouter(c context.Context, id uuid.UUID) (*model.Router, error) {
	return FromContext(c).GetRouter(id)
}

//GetRoutesList gets all active routes
func GetRoutesList(c context.Context) ([]*model.Router, error) {
	return FromContext(c).GetRoutesList()
}
