package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

// TODO: Switch migration package to github.com/mattes/migrate
func RunMigration(db *pg.DB, arg ...string) (int64, int64, error) {
	return migrations.Run(db, arg...)
}
