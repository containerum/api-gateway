package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//Router is struct for keep user proxy task
type Router struct {
	ID    uuid.UUID `sql:"type:uuid"`
	Group string

	Active  bool
	Created time.Time
}
