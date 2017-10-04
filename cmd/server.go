package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GATEWAY_DEBUG",
		Name:   "debug, d",
		Usage:  "start the server in debug mode",
	},
	cli.StringFlag{
		EnvVar: "PG_USER",
		Name:   "pg-user",
		Usage:  "Postgres user",
		Value:  "pg",
	},
	cli.StringFlag{
		EnvVar: "PG_PASSWORD",
		Name:   "pg-password",
		Usage:  "Postgres user password",
		Value:  "123456789",
	},
	cli.StringFlag{
		EnvVar: "PG_DATABASE",
		Name:   "pg-database",
		Usage:  "Postgres database",
		Value:  "postgres",
	},
	cli.StringFlag{
		EnvVar: "PG_ADDRESS",
		Name:   "pg-address",
		Usage:  "Postgres address",
		Value:  "x1.containerum.io:36519",
	},
}

const usageText = `Awesome Golang API Gateway.

	Migrations! This program runs command on the db. Supported commands are:
  - init - creates gopg_migrations table.
  - up - runs all available migrations.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
`

func server(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		log.Debug("Application running in Debug mode")
	} else {
		log.SetLevel(log.WarnLevel)
	}

	//Setup store
	log.WithFields(log.Fields{
		"PG_USER":     c.String("pg-user"),
		"PG_PASSWORD": c.String("pg-password"),
		"PG_DATABASE": c.String("pg-database"),
		"PG_ADDRESS":  c.String("pg-address"),
	}).Debug("Setup DB connection")
	setupStore(c)

	return nil
}
