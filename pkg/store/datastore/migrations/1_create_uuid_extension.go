package migrations

import (
	"github.com/jinzhu/gorm"
)

func init() {
	addMigration(Migration{
		Version: 1,
		up: func(db *gorm.DB) (int, error) {
			return 1, db.Raw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
		},
		down: func(db *gorm.DB) (int, error) {
			return 1, db.Raw(`DROP EXTENSION IF EXISTS "uuid-ossp";`).Error
		},
	})
}
