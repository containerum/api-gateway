package datastore

import (
	"bitbucket.org/exonch/ch-gateway/pkg/model"
	uuid "github.com/satori/go.uuid"
)

func (db *datastore) GetRouter(uuid.UUID) (*model.Router, error) {
	return nil, nil
}

func (db *datastore) GetRoutesList() ([]*model.Router, error) {
	return nil, nil
}
