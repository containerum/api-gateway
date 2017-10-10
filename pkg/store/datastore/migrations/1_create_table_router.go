package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table router")
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS router`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table router")
		_, err := db.Exec(`DROP TABLE router`)
		return err
	})
}
