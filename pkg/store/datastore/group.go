package datastore

import (
	"git.containerum.net/ch/api-gateway/pkg/model"
)

func (d *datastore) GetGroupList(g *model.Group) (*model.Group, error) {
	var group model.Group
	err := d.Where(g).Find(&group).Error
	return &group, err
}

func (d *datastore) CreateGroup(g *model.Group) (*model.Group, error) {
	d.NewRecord(g)
	return g, d.Save(g).Error
}
