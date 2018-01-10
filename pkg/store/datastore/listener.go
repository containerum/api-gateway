package datastore

import (
	"errors"
	"fmt"

	"git.containerum.net/ch/api-gateway/pkg/model"
)

//GetListener find Listener by ID
func (d *datastore) GetListener(id string) (*model.Listener, error) {
	var listener model.Listener
	d.Where("id = ?", id).First(&listener)
	if listener.ID == "" {
		return &listener, fmt.Errorf("Unable to find Listener with id = %v", id)
	}
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
func (d *datastore) UpdateListener(l *model.Listener, utype int) error {
	return errors.New("Not allowed update type")
}

//CreateListener create new listener in DB
func (d *datastore) CreateListener(l *model.Listener) (*model.Listener, error) {
	d.NewRecord(l)
	return l, d.Save(l).Error
}

//DeleteListener delete listener in DB by ID
func (d *datastore) DeleteListener(id string) error {
	return d.Where("id = ?", id).Delete(&model.Listener{}).Error
}
