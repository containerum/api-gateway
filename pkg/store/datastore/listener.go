package datastore

import (
	"bitbucket.org/exonch/ch-gateway/pkg/model"
)

//GetListener find Listener by ID
func (d *datastore) GetListener(id string) (*model.Listener, error) {
	var listener model.Listener

	return &listener, nil
}

//FindListener find first listener by input model
func (d *datastore) FindListener(l *model.Listener) (*model.Listener, error) {
	var listener model.Listener

	return &listener, nil
}

//GetListenerList find all listeers by input model
func (d *datastore) GetListenerList(l *model.Listener) (*[]model.Listener, error) {
	var listeners []model.Listener
	err := d.Where(l).Find(&listeners).Error
	return &listeners, err
}

//UpdateListener updates model in DB
func (d *datastore) UpdateListener(l *model.Listener) error {
	return d.Save(l).Error
}

//CreateListener create new listener in DB
func (d *datastore) CreateListener(l *model.Listener) (*model.Listener, error) {
	d.NewRecord(l)
	return l, d.Save(l).Error
}
