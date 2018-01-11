package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//Migration keeps migration currect migration version in PG
type Migration struct {
	ID      uint `gorm:"primary_key"`
	Version int  `gorm:"not null"`
	Dirty   bool `gorm:"not null"`

	up   func(*gorm.DB) (int, error)
	down func(*gorm.DB) (int, error)
}

var migrationList map[int]Migration

func addMigration(m Migration) {
	if migrationList == nil {
		migrationList = make(map[int]Migration)
	}
	migrationList[m.Version] = m
}

//RunMigrations run all migrations
func RunMigrations(d *gorm.DB) (int, error) {
	i, k := 0, 0
	for i <= len(migrationList) {
		if m, ok := migrationList[i]; ok {
			if err := m.startMigration(d); err != nil {
				return i, fmt.Errorf("Start migration failed. Err: %v", err)
			}
			if v, err := m.up(d); err != nil {
				return v, err
			}
			if err := m.finishMigration(d); err != nil {
				return i, fmt.Errorf("Finish migration failed. Err: %v", err)
			}
			k = i
		}
		i++
	}
	return k, nil
}

//RunMigrationsDown last down migration
func RunMigrationsDown(v int, d *gorm.DB) (int, error) {
	if v > 0 {
		if m, ok := migrationList[v]; ok {
			if err := m.startMigration(d); err != nil {
				return v, fmt.Errorf("Start migration failed. Err: %v", err)
			}
			if _, err := m.down(d); err != nil {
				return v, err
			}
			if err := m.finishMigration(d); err != nil {
				return v, fmt.Errorf("Finish migration failed. Err: %v", err)
			}
			if err := m.decreaseMigration(d); err != nil {
				return v, fmt.Errorf("Decrease migration failed. Err: %v", err)
			}
		} else {
			return v, fmt.Errorf("Unable to find migration %v", v)
		}
	}
	return 0, nil
}

func (m *Migration) startMigration(d *gorm.DB) error {
	var firstMigration Migration
	err := d.First(&firstMigration).Error
	if err != nil {
		return fmt.Errorf("Unable to get first migration record. Err: %v", err)
	}
	firstMigration.Version = m.Version
	firstMigration.Dirty = true
	return d.Model(&firstMigration).UpdateColumn("version", "dirty").Error
}

func (m *Migration) finishMigration(d *gorm.DB) error {
	var firstMigration Migration
	err := d.First(&firstMigration).Error
	if err != nil {
		return fmt.Errorf("Unable to get first migration record. Err: %v", err)
	}
	firstMigration.Version = m.Version
	firstMigration.Dirty = false
	return d.Model(&firstMigration).UpdateColumn("version", "dirty").Error
}

func (m *Migration) decreaseMigration(d *gorm.DB) error {
	var firstMigration Migration
	err := d.First(&firstMigration).Error
	if err != nil {
		return fmt.Errorf("Unable to get first migration record. Err: %v", err)
	}
	firstMigration.Version = m.Version - 1
	fmt.Print(firstMigration)
	return d.Model(&firstMigration).UpdateColumn("version", "dirty").Error
}
