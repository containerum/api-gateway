package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("Creating table: Routers")
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS routers (
				id uuid PRIMARY KEY,
				group_id uuid NOT NULL,
				roles text[],
				o_auth boolean DEFAULT true,
				active boolean DEFAULT true,
				created timestamp with time zone DEFAULT current_timestamp
			)
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("Dropping table: Routers")
		_, err := db.Exec(`DROP TABLE routers`)
		return err
	})
}
