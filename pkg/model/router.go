package model

import (
	"time"

	"github.com/satori/go.uuid"
)

//Router is struct for keep user proxy task
type Router struct {
	tableName struct{} `sql:"routers,alias:router"`

	ID     string    `sql:"type:uuid"`
	IDuuid uuid.UUID `sql:"-"`

	Group   *Group
	GroupID string

	Roles   []string
	OAuth   bool
	Active  bool
	Created time.Time
}

//BeforeInsert set GroupID
func (r *Router) BeforeInsert() {
	if r.GroupID == "" {
		r.GroupID = r.Group.ID
	}
	if r.ID == "" {
		r.ID = r.IDuuid.String()
	}
}
