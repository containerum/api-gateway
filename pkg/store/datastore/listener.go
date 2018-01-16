package datastore

import (
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

//GetListener find Listener by ID
func (d *data) GetListener(id string) (*model.Listener, error) {
	var listener model.Listener
	rows, err := d.Queryx(SQLGetListener, id)
	if err != nil {
		log.WithError(err).Warn(ErrUnableGetListener)
		return nil, ErrUnableGetListener
	}
	defer rows.Close()
	if !rows.Next() {
		log.Warn(ErrNoRows)
		return nil, ErrNoRows
	}
	if err := rows.StructScan(&listener); err != nil {
		log.WithError(err).Warn(ErrUnableScanListener)
		return nil, ErrUnableScanListener
	}
	return &listener, nil
}

//GetListenerList find all listeers by input model
func (d *data) GetListenerList(active *bool) (*[]model.Listener, error) {
	var listeners []model.Listener
	var err error
	rows := initRows()
	if active != nil {
		rows, err = d.Queryx(SQLGetListenersActive, active)
	} else {
		rows, err = d.Queryx(SQLGetListeners)
	}
	if err != nil {
		log.WithError(err).Warn(ErrUnableGetListeners)
		return nil, ErrUnableGetListeners
	}
	defer rows.Close()
	for rows.Next() {
		var listener model.Listener
		if err := rows.StructScan(&listener); err != nil {
			log.WithError(err).Warn(ErrUnableScanListener)
			return nil, ErrUnableScanListener
		}
		listeners = append(listeners, listener)
	}
	return &listeners, nil
}

//TODO Get updated time
//UpdateListener updates model in DB
func (d *data) UpdateListener(l *model.Listener) error {
	_, err := d.Exec(SQLUpdateListener, l.Name, l.OAuth, l.Active, l.StripPath, l.ListenPath, l.UpstreamURL, l.Method, l.GroupRefer, l.ID)
	if err != nil {
		log.WithError(err).Warn(ErrUnableUpdateListener)
		return ErrUnableUpdateListener
	}
	return nil
}

//CreateListener create new listener in DB
func (d *data) CreateListener(l *model.Listener) (*model.Listener, error) {
	var id string
	var created time.Time
	rows, err := d.Queryx(SQLCreateListener, l.Name, l.OAuth, l.Active, l.StripPath, l.ListenPath, l.UpstreamURL, l.Method, l.GroupRefer)
	if err != nil {
		log.WithError(err).Warn(ErrUnableCreateListener)
		return nil, ErrUnableCreateListener
	}
	defer rows.Close()
	if !rows.Next() {
		log.Warn(ErrNoRows)
		return nil, ErrNoRows
	}
	if err := rows.Scan(&id, &created); err != nil {
		log.WithError(err).Warn(ErrUnableScanListenerID)
		return nil, ErrUnableScanListenerID
	}
	l.ID = id
	l.CreatedAt, l.UpdatedAt = created, created
	return l, nil
}

//DeleteListener delete listener in DB by ID
func (d *data) DeleteListener(id string) error {
	_, err := d.Exec(SQLDeleteListener, id)
	if err != nil {
		log.WithError(err).Warn(ErrUnableDeleteListener)
		return ErrUnableDeleteListener
	}
	return nil
}
