package main

import (
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/router"
	"git.containerum.net/ch/api-gateway/pkg/router/middleware"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func runServer(c *cli.Context) error {
	setLogFormat(c) //Set log format
	log.WithFields(log.Fields{
		"PG_USER":       c.String("pg-user"),
		"PG_PASSWORD":   c.String("pg-password"),
		"PG_DATABASE":   c.String("pg-database"),
		"PG_ADDRESS":    c.String("pg-address"),
		"PG_PORT":       c.String("pg-port"),
		"PG_MIGRATIONS": c.Bool("pg-migrations"),
	}).Debug("Setup DB connection")

	//Create router and register all extensions
	r := router.CreateRouter(nil)
	r.RegisterStore(setupStore(c))
	r.RegisterAuth(setupAuth(c))
	r.RegisterRatelimiter(setupRatelimiter(c))
	r.RegisterStatsd(setupStatsd(c))
	r.RegisterClickhouseLogger(setupClickhouseLogger(c))
	r.InitRoutes()
	r.Start()

	return listenAndServe(r)
}

//GetVersion return app version
func GetVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}

func setLogFormat(c *cli.Context) error {
	if c.Bool("debug") {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
		log.Debug("Application running in Debug mode")
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	}
	return nil
}

func listenAndServe(handler http.Handler) error {
	//TODO: Move Cors to middleware
	c := middleware.Cors()
	server := &http.Server{Addr: ":8082", Handler: c.Handler(handler)}
	return server.ListenAndServe()
}
