package datastore

import (
	"bitbucket.org/exonch/ch-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

//GetRouter return router by UUID
func (db *datastore) GetRouter(id string) (*model.Router, error) {
	r := &model.Router{ID: id}
	err := db.Select(r)
	return r, err
}

//GetRoutesList return list of all routers
func (db *datastore) GetRoutesList() (*[]model.Router, error) {
	var rs []model.Router
	err := db.Model(&rs).Select()
	return &rs, err
}

//GetRoutesListByActivation return list of all routers
func (db *datastore) GetRoutesListByActivation(active bool) (*[]model.Router, error) {
	var rs []model.Router
	r := &model.Router{Active: active}
	err := db.Model(&rs).Select(r)
	return &rs, err
}

//AddRouter create new router
func (db *datastore) AddRouter(r *model.Router) error {
	r.BeforeInsert()
	if err := db.Insert(r); err != nil {
		log.WithFields(log.Fields{
			"Router": string(r.ConvertToJSON()),
			"Err":    err,
		}).Error("AddRouter error")
		return err
	}
	log.WithFields(log.Fields{
		"Router": string(r.ConvertToJSON()),
	}).Debug("AddRouter success")
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
