package main

import (
	"fmt"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"github.com/BurntSushi/toml"
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

	// sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	// go func() {
	// 	sig := <-sigs
	// 	fmt.Println(sig)
	// 	done <- true
	// }()
	//
	// fmt.Println("awaiting signal")
	// <-done
	// fmt.Println("exiting")
	//
	// return nil

	// setupLogger(c)
	// serve := setupServer(c)
	// if conf.TLS.Enable {
	// 	cert, key, err := setupTSL(c)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return serve.Run(port, cert, key)
	// }
	// return serve.Run(port, "", "")

	var routes model.Routes
	if _, err := toml.DecodeFile("routes.toml", &routes); err != nil {
		return err
	}
	fmt.Println(routes)

	return nil
}

//GetVersion return app version
func GetVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}
