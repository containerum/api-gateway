package model

import (
	"time"
)

//DatabaseConfig containts DB config
type DatabaseConfig struct {
	User          string
	Password      string
	Database      string
	Address       string
	Port          string
	SafeMigration bool
	Debug         bool
	Migrations    bool
}

const timeFormat = "Monday, 02 Jan 2006 15:04:05 MST"

//DefaultModel embedded struct in each DB Model
type DefaultModel struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func getUnixTimeOrNil(t *time.Time) *int64 {
	var unix int64
	if t != nil {
		unix = t.Unix()
		return &unix
	}
	return nil
}
