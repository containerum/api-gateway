package model

import uuid "github.com/satori/go.uuid"

//Group struct ident router group
type Group struct {
	tableName struct{} `sql:"groups"`

	ID     string `sql:"type:uuid" json:"-"`
	Name   string `sql:"type:varchar(32)"`
	Active bool
}

//CreateDefaultGroup return default Group
func CreateDefaultGroup() *Router {
	return &Router{
		ID:     uuid.NewV4().String(),
		Active: true,
	}
}
