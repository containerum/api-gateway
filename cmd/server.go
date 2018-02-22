package main

import (
	"github.com/urfave/cli"
)

var (
	conf Config
	port = 8082
)

func runServer(c *cli.Context) error {
	var err error
	if conf, err = getConfig(); err != nil {
		return err
	}
	setupLogger(c)
	serve := setupServer(c)
	if conf.TLS.Enable {
		cert, key, err := setupTSL(c)
		if err != nil {
			return err
		}
		return serve.Run(port, cert, key)
	}
	return serve.Run(port, "", "")
}

//GetVersion return app version
func GetVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}
