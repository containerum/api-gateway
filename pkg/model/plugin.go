// +build ignore

package model

//Plugin keeps plagin name and state when it should be work: After proxy or before
type Plugin struct {
	ID         uint   `gorm:"primary_key"`
	ListenerID string `gorm:"type:uuid;not null"`
	Name       string `gorm:"not null"`
	RunAfter   bool   `gorm:"not null;default:true"`
}
