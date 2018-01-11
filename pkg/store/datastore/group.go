package datastore

import (
	"errors"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

func (d *datastore) GetGroupList(g *model.Group) (*[]model.Group, error) {
	reqName := "GetGroupList call"
	var groups []model.Group
	err := d.Where(g).Find(&groups).Error
	if err != nil {
		log.WithError(err).Error(reqName)
		return nil, errors.New("")
	}
	log.Debug(reqName)
	return &groups, nil
}

func (d *datastore) CreateGroup(g *model.Group) (*model.Group, error) {
	d.NewRecord(g)
	return g, d.Save(g).Error
}
