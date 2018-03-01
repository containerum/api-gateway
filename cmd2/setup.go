package main

import (
	"errors"
	"fmt"
	"os"

	"git.containerum.net/ch/api-gateway/pkg2/model"
	"git.containerum.net/ch/api-gateway/pkg2/server"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type tomlFile interface {
	Validate() []error
}

var (
	config     model.Config
	routes     model.Routes
	cert, key  string
	authClient *auth.AuthClient
)

var (
	errUnableReadConfig   = errors.New("Unable to read config file")
	errUnableReadRoutes   = errors.New("Unable to read routes file")
	errUnableOpenCertFile = errors.New("Unable to open cert.pem")
	errUnableOpenKeyFile  = errors.New("Unable to open key.pem")

	errGrpcDialFailed = errors.New("Dial to auth grpc failed")
)

//TODO: parse flags
func setupLogs(c *cli.Context) {
	gin.SetMode(gin.ReleaseMode)
	log.SetLevel(log.DebugLevel)
}

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
		log.WithError(e).Error(errUnableOpenCertFile)
		return errUnableOpenCertFile
	}
	if _, e := os.Stat(c.String(tlsKeyPath)); os.IsNotExist(e) {
		log.WithError(e).Error(errUnableOpenKeyFile)
		return errUnableOpenKeyFile
	}
	cert = c.String(tlsCertPath)
	key = c.String(tlsKeyPath)
	config.TLS.Cert = cert
	config.TLS.Key = key
	log.WithField("Key", key).WithField("Cert", cert).Debug("TLS")
	return nil
}

func setupAuth(c *cli.Context) error {
	if !config.Auth.Enable {
		return nil
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	addr := fmt.Sprintf("%s:%s", c.String(authAddr), c.String(authPort))
	con, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error(errGrpcDialFailed)
		return errGrpcDialFailed
	}
	client := auth.NewAuthClient(con)
	authClient = &client
	return nil
}

func setupServer(c *cli.Context) (*server.Server, error) {
	return server.New(&config, &routes, authClient)
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
