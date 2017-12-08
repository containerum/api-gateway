package model

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel `structs:"-"`
	Name         string `gorm:"not null" structs:"name"`

	Roles []Role `structs:"roles"`
	//TODO: Add Group
	OAuth  bool `structs:"o_auth"`
	Active bool `structs:"active"`

	StripPath   bool   `structs:"strip_path"`
	ListenPath  string `gorm:"not null" structs:"listen_path"`
	UpstreamURL string `gorm:"not null" structs:"upstream_url"`
	Method      string `gorm:"not null" structs:"method"`

	Plugins []Plugin `structs:"plugins"`
}

func (l Listener) MarshalJSON() ([]byte, error) {
	listenerMap := structs.Map(l)
	for k, v := range l.getMap() {
		listenerMap[k] = v
	}
	return json.Marshal(listenerMap)
}

func (lm *Listener) BeforeUpdate() (err error) {
	fmt.Printf("\nUpdate: %v\n", lm.UpdatedAt)
	return
}
