package migrations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Info("Creating table: Routers")
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS routers (
				id uuid PRIMARY KEY,
				name varchar(32) NOT NULL,
				group_id uuid NOT NULL,
				roles text[],
				o_auth boolean,
				active boolean,
				created timestamp with time zone DEFAULT current_timestamp,
				strip_path boolean,
				listen_path varchar(128) NOT NULL,
				upstream_url varchar(256) NOT NULL,
				methods text[],
				before_plugins text[],
				after_plugins text[]
			)
		`)
		return err
	}, func(db migrations.DB) error {
		log.Info("Dropping table: Routers")
		_, err := db.Exec(`DROP TABLE routers`)
		return err
	})
}
