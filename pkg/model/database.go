package model

import (
	"bytes"
	"fmt"
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
}

//JSONTime return special format in output JSON
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Monday, 02 Jan 2006 15:04:05 MST"))
	return []byte(stamp), nil
}

//DefaultModel embedded struct in each Gorm Model
type DefaultModel struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (dm *DefaultModel) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	// buffer.WriteString(fmt.Sprintf("\"%s\":%s", "id", dm.ID))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
