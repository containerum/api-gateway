package main

import (
	"net/http"
	"time"

	"bitbucket.org/exonch/ch-gateway/pkg/router"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GATEWAY_DEBUG",
		Name:   "debug, d",
		Usage:  "start the server in debug mode",
	},
	cli.BoolFlag{
		Name:  "migrate, m",
		Usage: "start the server in migration mode",
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

	 Migrations runs only migrate mode! Supported commands are:
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
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}

	//Setup store
	log.WithFields(log.Fields{
		"PG_USER":     c.String("pg-user"),
		"PG_PASSWORD": c.String("pg-password"),
		"PG_DATABASE": c.String("pg-database"),
		"PG_ADDRESS":  c.String("pg-address"),
	}).Debug("Setup DB connection")
	// setupStore(c)

	r := router.CreateRouter()

	go func(r *router.Router) {
		time.Sleep(time.Second * 10)
		r.AddRoute("x1")
	}(r)

	go func(r *router.Router) {
		time.Sleep(time.Second * 20)
		r.AddRoute("x2")
	}(r)

	go func(r *router.Router) {
		time.Sleep(time.Second * 30)
		r.AddRoute("x3")
	}(r)

	return listenAndServe(r)
}

func listenAndServe(handler http.Handler) error {
	server := &http.Server{Addr: ":8081", Handler: handler}
	return server.ListenAndServe()
}
