// +build ignore

package model

//Role keeps user role
type Role struct {
	ID         uint   `gorm:"primary_key"`
	ListenerID string `gorm:"type:uuid;not null"`
	Name       string `gorm:"not null"`
}
