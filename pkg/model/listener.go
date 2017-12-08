package model

import "fmt"

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel
	Name string `gorm:"not null"`

	Roles []Role
	//TODO: Add Group
	OAuth  bool
	Active bool

	StripPath   bool
	ListenPath  string `gorm:"not null"`
	UpstreamURL string `gorm:"not null"`
	Method      string `gorm:"not null"`

	Plugins []Plugin
}

func (lm *Listener) BeforeUpdate() (err error) {
	fmt.Printf("\nUpdate: %v\n", lm.UpdatedAt)
	return
}
