package main

import (
	"errors"
	"fmt"
	"os"

	"git.containerum.net/ch/api-gateway/pkg/gatewayErrors"
	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server"
	toml "git.containerum.net/ch/api-gateway/pkg/utils/toml"
	"git.containerum.net/ch/auth/proto"
	"github.com/containerum/cherry/adaptors/cherrygrpc"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"

	log "github.com/sirupsen/logrus"
)

type setupFunc func(c *cli.Context) error

const (
	configKey     = "config"
	routesKey     = "routes"
	authClientKey = "auth"
	metricsKey    = "metrics"
)

var (
	errUnableReadConfig   = errors.New("unable to read config file")
	errUnableReadRoutes   = errors.New("unable to read routes file")
	errUnableOpenCertFile = errors.New("unable to open cert.pem")
	errUnableOpenKeyFile  = errors.New("unable to open key.pem")

	errGrpcDialFailed = errors.New("dial to auth grpc failed")
)

func setup(c *cli.Context, fns ...setupFunc) (err error) {
	for _, fn := range fns {
		if err := fn(c); err != nil {
			return err
		}
	}
	return
}

func setupLogs(c *cli.Context) (err error) {
	gin.SetMode(gin.ReleaseMode)
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
	return
}

func setupConfig(c *cli.Context) (err error) {
	cfg := &model.Config{}
	if err := toml.ReadToml(c.String(ConfigPathFlag.Name), cfg); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadConfig, err)
	}
	log.Infof("Config setup:\n%v", cfg)
	c.App.Metadata[configKey] = cfg
	return
}

func setupRoutes(c *cli.Context) (err error) {
	routes := &model.Routes{}
	if err = toml.ReadToml(c.String(RoutesPathFlag.Name), routes); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadRoutes, err)
	}
	for key := range routes.Routes {
		if route, ok := routes.Routes[key]; ok {
			route.ID = key
			routes.Routes[key] = route
		}
	}
	c.App.Metadata[routesKey] = routes
	return
}

func setupTLS(c *cli.Context) (err error) {
	config := c.App.Metadata[configKey].(*model.Config)
	if !config.TLS.Enable {
		return
	}
	cert, key := c.String(TLSCertPathFlag.Name), c.String(TLSKeyPathFlag.Name)
	if _, e := os.Stat(cert); os.IsNotExist(e) {
		log.WithError(e).Error(errUnableOpenCertFile)
		return errUnableOpenCertFile
	}
	if _, e := os.Stat(key); os.IsNotExist(e) {
		log.WithError(e).Error(errUnableOpenKeyFile)
		return errUnableOpenKeyFile
	}
	config.TLS.Cert, config.TLS.Key = cert, key
	log.WithField("Key", key).WithField("Cert", cert).Debug("TLS")
	return
}

func setupAuth(c *cli.Context) (err error) {
	config := c.App.Metadata[configKey].(*model.Config)
	if !config.Auth.Enable {
		c.App.Metadata[authClientKey] = (authProto.AuthClient)(nil)
		return
	}
	opts := append([]grpc.DialOption{}, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		grpc_logrus.UnaryClientInterceptor(log.WithField("component", "auth_client")),
		cherrygrpc.UnaryClientInterceptor(gatewayErrors.ErrInternal),
	)))
	var con *grpc.ClientConn
	if con, err = grpc.Dial(c.String(AuthAddrFlag.Name), opts...); err != nil {
		log.WithError(err).Error(errGrpcDialFailed)
		return errGrpcDialFailed
	}
	c.App.Metadata[authClientKey] = authProto.NewAuthClient(con)
	return
}

func setupServer(c *cli.Context) (*server.Server, error) {
	opt := &server.Options{
		Routes:  c.App.Metadata[routesKey].(*model.Routes),
		Config:  c.App.Metadata[configKey].(*model.Config),
		Auth:    c.App.Metadata[authClientKey].(authProto.AuthClient),
		Metrics: c.App.Metadata[metricsKey].(*model.Metrics),

		ServiceHostPrefix: c.String(ServiceHostPrefixFlag.Name),
		Version:           c.App.Version,
	}
	return server.New(opt)
}

func setupMetrics(c *cli.Context) (err error) {
	config := c.App.Metadata[configKey].(*model.Config)
	metrics := model.CreateMetrics()
	if config.Prometheus.Enable {
		prometheus.MustRegister(metrics.RTotal, metrics.RUserIP, metrics.RRoute, metrics.RUserAgent)
	}
	c.App.Metadata[metricsKey] = metrics
	return
}
