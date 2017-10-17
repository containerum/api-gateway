package main

import (
	"net/http"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
	"bitbucket.org/exonch/ch-gateway/pkg/router"
	"bitbucket.org/exonch/ch-gateway/pkg/router/middleware"

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
	cli.StringFlag{
		EnvVar: "STATSD_ADDRESS",
		Name:   "statsd-address",
		Usage:  "Statsd address",
		Value:  "213.239.208.25:8125",
	},
	cli.StringFlag{
		Name:  "statsd-prefix",
		Usage: "Statsd prefix",
		Value: "ch-gateway",
	},
	cli.IntFlag{
		Name:  "statsd-buffer-time",
		Usage: "Statsd buffer time",
		Value: 300,
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

	//Write store logs
	log.WithFields(log.Fields{
		"PG_USER":     c.String("pg-user"),
		"PG_PASSWORD": c.String("pg-password"),
		"PG_DATABASE": c.String("pg-database"),
		"PG_ADDRESS":  c.String("pg-address"),
	}).Debug("Setup DB connection")

	//Setup store
	s := setupStore(c)
	std := setupStatsd(c)

	//Create routers
	r := router.CreateRouter(&s, &std)

	m := model.CreateDefaultRouter()
	m.UpstreamURL = "https://web.containerum.io"
	m.ListenPath = "/xx"
	m.Methods = []string{"get"}
	r.AddRoute(m)

	return listenAndServe(r)
}

func listenAndServe(handler http.Handler) error {
	c := middleware.Cors()
	server := &http.Server{Addr: ":8080", Handler: c.Handler(handler)}
	return server.ListenAndServe()
}
