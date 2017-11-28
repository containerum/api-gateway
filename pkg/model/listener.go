package model

type Listener struct {
	DefaultModel
	Name string `gorm:"not null"`

	Roles  []Role
	OAuth  bool
	Active bool

	StripPath   bool
	ListenPath  string `gorm:"not null"`
	UpstreamURL string `gorm:"not null"`
	//Methods     []string `gorm:"type:varchar[]"`

	//BeforePlugins []string `gorm:"type:varchar[]"`
	//AfterPlugins  []string `gorm:"type:varchar[]"`

	//TODO: Add Group
}

type Role struct {
	ID         uint   `gorm:"primary_key"`
	ListenerID string `gorm:"type:uuid;not null"`
	Name       string `gorm:"not null"`
}
