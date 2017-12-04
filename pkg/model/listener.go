package model

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
