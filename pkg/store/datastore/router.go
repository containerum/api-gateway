package datastore

import (
	"bitbucket.org/exonch/ch-gateway/pkg/model"
	uuid "github.com/satori/go.uuid"
)

//GetRouter return router by UUID
func (db *datastore) GetRouter(uuid.UUID) (*model.Router, error) {
	return nil, nil
}

//GetRoutesList return list of all routers
func (db *datastore) GetRoutesList() ([]*model.Router, error) {
	return nil, nil
}

//GetRoutesListActivation return list of all routers
func (db *datastore) GetRoutesListByActivation(active bool) ([]*model.Router, error) {
	return nil, nil
}

//AddRouter create new router
func (db *datastore) AddRouter(r *model.Router) error {
	return nil
}

//RemoveRouter delete router from DB
func (db *datastore) RemoveRouter(r *model.Router) error {
	return nil
}

//UpdateRouter update router
func (db *datastore) UpdateRouter(r *model.Router) error {
	return nil
}

//SetRouterActivity change router active
func (db *datastore) SetRouterActivity(r *model.Router, active bool) error {
	return nil
}
