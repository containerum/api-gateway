// +build !stub

package datastore

import (
	"strings"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

//GetListener find Listener by ID
func (d *datastore) GetListener(id string) (*model.Listener, error) {
	reqName := "GetListener call"
	var listener model.Listener
	d.Where("id = ?", id).First(&listener)
	if listener.ID == "" {
		log.WithError(ErrListenerIDIsEmpty).Error(reqName)
		return &listener, ErrUnableFindListener
	}
	log.Debug(reqName)
	return &listener, nil
}

//TODO: Remove error from answer
//GetListenerList find all listeers by input model
func (d *datastore) GetListenerList(l *model.Listener) (*[]model.Listener, error) {
	reqName := "GetListenerList call"
	var listeners []model.Listener
	err := d.Where(l).Find(&listeners).Error
	if err != nil {
		log.WithError(err).Error(reqName)
	}
	log.Debug(reqName)
	return &listeners, nil
}

//UpdateListener updates model in DB
func (d *datastore) UpdateListener(l *model.Listener, utype model.ListenerUpdateType) error {
	reqName := "UpdateListener call"
	if _, err := d.GetListener(l.ID); err != nil {
		log.WithError(err).Error(reqName)
		return err
	}
	switch utype {
	case model.ListenerUpdateActive:
		err := d.Model(l).Update("active", l.Active).Error
		if err != nil {
			log.WithError(err).Error(reqName)
			return ErrUnableToUpdateListener
		}
		log.WithField("UpdateType", "ListenerUpdateActive").Debug(reqName)
	case model.ListenerUpdateOAuth:
		err := d.Model(l).Update("o_auth", l.OAuth).Error
		if err != nil {
			log.WithError(err).Error(reqName)
			return ErrUnableToUpdateListener
		}
		log.WithField("UpdateType", "ListenerUpdateOAuth").Debug(reqName)
	case model.ListenerUpdateFull:
		err := d.Model(l).Update(
			"name", l.Name,
			"method", strings.ToUpper(l.Method),
			"group_refer", l.GroupRefer,
			"listen_path", l.ListenPath,
			"upstream_url", l.UpstreamURL,
		).Error
		if err != nil {
			log.WithError(err).Error(reqName)
			return ErrUnableToUpdateListener
		}
		log.WithField("UpdateType", "ListenerUpdateFull").Debug(reqName)
	}
	return nil
}

//CreateListener create new listener in DB
func (d *datastore) CreateListener(l *model.Listener) (*model.Listener, error) {
	reqName := "CreateListener call"
	d.NewRecord(l)
	err := d.Save(l).Error
	if err != nil {
		log.WithError(err).Error(reqName)
		return nil, ErrUnableToCreateListener
	}
	log.Debug(reqName)
	return l, nil
}

//DeleteListener delete listener in DB by ID
func (d *datastore) DeleteListener(id string) error {
	reqName := "DeleteListener call"
	err := d.Where("id = ?", id).Delete(&model.Listener{}).Error
	if err != nil {
		log.WithError(err).Error(reqName)
		return ErrUnableToDeleteListener
	}
	log.Debug(reqName)
	return nil
}
