package main

import (
	"net/http"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/router"
	"git.containerum.net/ch/api-gateway/pkg/router/middleware"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func runServer(c *cli.Context) error {
	setLogFormat(c) //Set log format
	log.WithFields(log.Fields{
		"PG_USER":           c.String("pg-user"),
		"PG_PASSWORD":       c.String("pg-password"),
		"PG_DATABASE":       c.String("pg-database"),
		"PG_ADDRESS":        c.String("pg-address"),
		"PG_PORT":           c.String("pg-port"),
		"PG_DEBUG":          c.Bool("pg-debug"),
		"PG_SAFE_MIGRATION": c.Bool("pg-safe-migration"),
	}).Debug("Setup DB connection") //Write store logs

	//Create router and register all extensions
	r := router.CreateRouter()
	r.RegisterStore(setupStore(c))
	r.RegisterAuth(setupAuth(c))
	r.RegisterRatelimiter(setupRatelimiter(c))
	r.RegisterStatsd(setupStatsd(c))
	r.RegisterClickhouseLogger(setupClickhouseLogger(c))
	r.InitRoutes()
	r.Start(time.Second * 5)

	return listenAndServe(r)
}

//Initialize migration table
func initMigration(c *cli.Context) error {
	s := *setupStore(c)
	if err := s.Init(); err != nil {
		log.WithField("Error", err.Error()).Error("Migration table creation is failed")
	}
	log.Info("Migration table is successfully created")
	return nil
}

//Print migration version
func getMigrationVersion(c *cli.Context) error {
	s := *setupStore(c)
	if v, err := s.Version(); err != nil {
		log.WithField("Error", err.Error()).Error("Unable to get migration version")
	} else {
		log.Infof("Migration version is %v", v)
	}
	return nil
}

//Run all migrations
func upMigration(c *cli.Context) error {
	s := *setupStore(c)
	if v, err := s.Up(); err != nil {
		log.WithField("Error", err.Error()).Error("Migration failed")
	} else {
		log.Infof("Migration is Up, Version is: %v", v)
	}
	return nil
}

//Run curent version down migration
func downMigration(c *cli.Context) error {
	s := *setupStore(c)
	if v, err := s.Down(); err != nil {
		log.WithField("Error", err.Error()).Error("Migration failed")
	} else {
		log.Infof("Migration is Down, Version is: %v", v)
	}
	return nil
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
	server := &http.Server{Addr: ":8081", Handler: c.Handler(handler)}
	return server.ListenAndServe()
}
