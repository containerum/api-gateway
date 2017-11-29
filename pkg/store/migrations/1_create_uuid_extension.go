package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Migration keeps migration currect migration version in PG
type Migration struct {
	ID      uint   `gorm:"primary_key"`
	Version int    `gorm:"not null"`
	Dirty   string `gorm:"type:uuid;not null"`

	Up   func(*gorm.DB) error `gorm:"-"`
	Down func(*gorm.DB) error `gorm:"-"`
}

func init() {
	addMigration(Migration{
		Version: 1,
		Up: func(db *gorm.DB) error {
			fmt.Print("UP")
			return nil
		},
	})
}
