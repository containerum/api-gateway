package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("Creating table: Groups")
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS groups (
				id uuid PRIMARY KEY,
        name varchar(16) UNIQUE,
				active boolean DEFAULT true
			)
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("Dropping table: Groups")
		_, err := db.Exec(`DROP TABLE groups`)
		return err
	})
}
