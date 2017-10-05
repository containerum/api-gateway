package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

func RunMigration(db *pg.DB, arg ...string) (int64, int64, error) {
	return migrations.Run(db, arg...)
}
