package model

import "github.com/satori/go.uuid"

//Group struct ident router group
type Group struct {
	tableName struct{} `sql:"groups"`

	ID     string    `sql:"type:uuid"`
	IDuuid uuid.UUID `sql:"-"`

	Name   string `sql:"type:varchar(32)"`
	Active bool
}

//BeforeInsert set ID
func (g *Group) BeforeInsert() {
	if g.ID == "" {
		g.ID = g.IDuuid.String()
	}
}
