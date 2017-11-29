package model

//Group struct ident router group
type Group struct {
	DefaultModel
	Name   string `gorm:"type:varchar(32);not null;unique"`
	Active bool   `gorm:"not null"`
}
