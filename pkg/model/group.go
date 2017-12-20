package model

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

//Group struct ident router group
type Group struct {
	DefaultModel `structs:"-" json:"-"`
	Name         string `gorm:"type:varchar(32);not null;unique" structs:"name"`
	Active       bool   `gorm:"not null" structs:"active"`
}

func (g Group) MarshalJSON() ([]byte, error) {
	groupMap := structs.Map(g)
	for k, v := range g.getMap() {
		groupMap[k] = v
	}
	return json.Marshal(groupMap)
}

func (g *Group) Valid() (err error) {
	if len(g.Name) < minNameLenght {
		return fmt.Errorf("Group name is too short. Min lenght: %d", minNameLenght)
	}
	return
}
