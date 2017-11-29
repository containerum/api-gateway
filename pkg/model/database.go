package model

import "time"

//DatabaseConfig containts DB config
type DatabaseConfig struct {
	User          string
	Password      string
	Database      string
	Address       string
	Port          string
	SafeMigration bool
	Debug         bool
}

//DefaultModel embedded struct in each Gorm Model
type DefaultModel struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
