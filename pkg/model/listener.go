package model

import "github.com/lib/pq"

//Listener keeps proxy-router configs
type Listener struct {
	DefaultModel
	Name string `gorm:"not null"`

	Roles []Role
	//TODO: Add Group
	OAuth  bool
	Active bool

	StripPath   bool
	ListenPath  string         `gorm:"not null"`
	UpstreamURL string         `gorm:"not null"`
	Methods     pq.StringArray `gorm:"type:varchar(16)[]"`

	Plugins []Plugin
}
