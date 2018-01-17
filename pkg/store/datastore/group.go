package datastore

import (
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"

	log "github.com/Sirupsen/logrus"
)

func (d *data) GetGroupList(g *model.Group) (*[]model.Group, error) {
	var groups []model.Group
	rows, err := d.Queryx(SQLGetGroups)
	if err != nil {
		log.WithError(err).Warn(ErrUnableGetGroups)
		return nil, ErrUnableGetGroups
	}
	defer rows.Close()
	for rows.Next() {
		var group model.Group
		if err := rows.StructScan(&group); err != nil {
			log.WithError(err).Warn(ErrUnableScanGroup)
			return nil, ErrUnableScanGroup
		}
		groups = append(groups, group)
	}
	return &groups, nil
}

func (d *data) CreateGroup(g *model.Group) (*model.Group, error) {
	var id string
	var created time.Time
	rows, err := d.Queryx(SQLCreateGroup, g.Name, g.Active)
	if err != nil {
		log.WithError(err).Warn(ErrUnableCreateGroup)
		return nil, ErrUnableCreateGroup
	}
	defer rows.Close()
	if !rows.Next() {
		log.Warn(ErrNoRows)
		return nil, ErrNoRows
	}
	if err := rows.Scan(&id, &created); err != nil {
		log.WithError(err).Warn(ErrUnableScanGroup)
		return nil, ErrUnableScanGroupID
	}
	g.ID = id
	g.CreatedAt, g.UpdatedAt = created, created
	return g, nil
}
