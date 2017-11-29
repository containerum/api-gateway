package migrations

import (
	"github.com/jinzhu/gorm"
)

var migrationList map[int]Migration

func addMigration(m Migration) {
	if migrationList == nil {
		migrationList = make(map[int]Migration)
	}
	migrationList[m.Version] = m
}

//RunMigrations run all migrations
func RunMigrations(d *gorm.DB) (int, error) {
	i := 1
	for ; i <= len(migrationList); i++ {
		if err := migrationList[i].Up(d); err != nil {
			return i, err
		}
	}
	return i, nil
}
