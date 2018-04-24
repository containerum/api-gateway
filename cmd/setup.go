package main

import (
	"errors"
	"fmt"
	"os"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server"
	toml "git.containerum.net/ch/api-gateway/pkg/utils/toml"
	"git.containerum.net/ch/auth/proto"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/cherrygrpc"
	"git.containerum.net/ch/kube-client/pkg/cherry/api-gateway"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type setupFunc func(c *cli.Context) error

var (
	config     model.Config
	routes     model.Routes
	cert, key  string
	authClient *authProto.AuthClient
	metrics    *model.Metrics
)

var (
	errUnableReadConfig   = errors.New("Unable to read config file")
	errUnableReadRoutes   = errors.New("Unable to read routes file")
	errUnableOpenCertFile = errors.New("Unable to open cert.pem")
	errUnableOpenKeyFile  = errors.New("Unable to open key.pem")

	errGrpcDialFailed = errors.New("Dial to auth grpc failed")
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
	if err := toml.ReadToml(c.String(configPath), &config); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadConfig, err)
	}
	return
}

func setupRoutes(c *cli.Context) (err error) {
	if err = toml.ReadToml(c.String(routesPath), &routes); err != nil {
		return fmt.Errorf("%v. %v", errUnableReadRoutes, err)
	}
	for key := range routes.Routes {
		if route, ok := routes.Routes[key]; ok {
			route.ID = key
			routes.Routes[key] = route
		}
	}
	return
}

func setupTLS(c *cli.Context) (err error) {
	if !config.TLS.Enable {
		return
	}
	if _, e := os.Stat(c.String(tlsCertPath)); os.IsNotExist(e) {
		log.WithError(e).Error(errUnableOpenCertFile)
		return errUnableOpenCertFile
	}
	if _, e := os.Stat(c.String(tlsKeyPath)); os.IsNotExist(e) {
		log.WithError(e).Error(errUnableOpenKeyFile)
		return errUnableOpenKeyFile
	}
	cert, key = c.String(tlsCertPath), c.String(tlsKeyPath)
	config.TLS.Cert, config.TLS.Key = cert, key
	log.WithField("Key", key).WithField("Cert", cert).Debug("TLS")
	return
}

func setupAuth(c *cli.Context) (err error) {
	if !config.Auth.Enable {
		return
	}
	opts := append([]grpc.DialOption{}, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		cherrygrpc.UnaryClientInterceptor(gatewayErrors.ErrInternal),
		grpc_logrus.UnaryClientInterceptor(log.WithField("component", "auth_client")),
	)))
	var con *grpc.ClientConn
	addr := fmt.Sprintf("%s:%s", c.String(authAddr), c.String(authPort))
	if con, err = grpc.Dial(addr, opts...); err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error(errGrpcDialFailed)
		return errGrpcDialFailed
	}
	client := authProto.NewAuthClient(con)
	authClient = &client
	return
}

func setupServer(c *cli.Context) (*server.Server, error) {
	opt := &server.Options{
		Routes:  &routes,
		Config:  &config,
		Auth:    authClient,
		Metrics: metrics,
	}
	return server.New(opt)
}

func setupMetrics(c *cli.Context) (err error) {
	if config.Prometheus.Enable {
		metrics = model.CreateMetrics()
		prometheus.MustRegister(metrics.RTotal, metrics.RUserIP, metrics.RRoute, metrics.RUserAgent)
	}
	return
}
