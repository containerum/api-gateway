package main

import (
	"errors"
	"fmt"
	"os"

	"git.containerum.net/ch/api-gateway/pkg2/model"
	"git.containerum.net/ch/api-gateway/pkg2/server"
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

type tomlFile interface {
	Validate() []error
}

var (
	config    model.Config
	routes    model.Routes
	cert, key string
)

var (
	errUnableReadConfig   = errors.New("Unable to read config file")
	errUnableReadRoutes   = errors.New("Unable to read routes file")
	errUnableOpenCertFile = errors.New("Unable to open cert.pem")
	errUnableOpenKeyFile  = errors.New("Unable to open key.pem")
)

func setupConfig(c *cli.Context) error {
	if err := readToml(c.String(configPath), &config); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadConfig, err)
	}
	return nil
}

func setupRoutes(c *cli.Context) error {
	if err := readToml(c.String(routesPath), &routes); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadRoutes, err)
	}
	return nil
}

func setupTLS(c *cli.Context) error {
	if !config.TLS.Enable {
		return nil
	}
	if _, e := os.Stat(c.String(tlsCertPath)); os.IsNotExist(e) {
		cert = c.String(tlsCertPath)
		config.TLS.Cert = cert
		log.WithError(e).Error(errUnableOpenCertFile)
		return errUnableOpenCertFile
	}
	if _, e := os.Stat(c.String(tlsKeyPath)); os.IsNotExist(e) {
		key = c.String(tlsKeyPath)
		config.TLS.Key = key
		log.WithError(e).Error(errUnableOpenKeyFile)
		return errUnableOpenKeyFile
	}
	return nil
}

func setupServer(c *cli.Context) (*server.Server, error) {

	log.SetLevel(log.DebugLevel)

	return server.New(&config, &routes)
}

func getVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}

func readToml(file string, out tomlFile) error {
	if _, err := toml.DecodeFile(file, out); err != nil {
		return err
	}
	if errs := out.Validate(); errs != nil {
		var err error
		for i, e := range errs {
			if i == 0 {
				err = fmt.Errorf("%s", e.Error())
			} else {
				err = fmt.Errorf("%s\n%s", err.Error(), e.Error())
			}
		}
		return err
	}
	return nil
}
