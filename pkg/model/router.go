package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//Router is struct for keep user proxy task
type Router struct {
	ID uuid.UUID `sql:"type:uuid"`

	Group *Group
	Roles []uint8

	OAuth  bool
	Active bool

	Created time.Time
}

//Group struct ident router group
type Group struct {
	ID   int
	Name string `sql:"type:varchar(32)"`
}
