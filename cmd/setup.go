package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server"
	"git.containerum.net/ch/api-gateway/pkg/utils/preproc"
	"git.containerum.net/ch/auth/proto"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/cherrygrpc"
	"git.containerum.net/ch/kube-client/pkg/cherry/api-gateway"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type tomlFile interface {
	Validate() []error
}

var (
	g errgroup.Group

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

//TODO: parse flags
func setupLogs(c *cli.Context) {
	gin.SetMode(gin.ReleaseMode)
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
		return
	}
	log.SetFormatter(&log.JSONFormatter{})
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
	for key := range routes.Routes {
		if route, ok := routes.Routes[key]; ok {
			route.ID = key
			routes.Routes[key] = route
		}
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

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		cherrygrpc.UnaryClientInterceptor(gatewayErrors.ErrInternal),
		grpc_logrus.UnaryClientInterceptor(log.WithField("component", "auth_client")),
	)))

	addr := fmt.Sprintf("%s:%s", c.String(authAddr), c.String(authPort))
	con, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.WithFields(log.Fields{
			"Err": err,
		}).Error(errGrpcDialFailed)
		return errGrpcDialFailed
	}
	client := authProto.NewAuthClient(con)
	authClient = &client
	return nil
}

func setupServer(c *cli.Context) (*server.Server, error) {
	opt := &server.ServerOptions{
		Routes:  &routes,
		Config:  &config,
		Auth:    authClient,
		Metrics: metrics,
	}
	return server.New(opt)
}

func setupMetrics(c *cli.Context) error {
	metrics = model.CreateMetrics()
	prometheus.MustRegister(metrics.RTotal, metrics.RUserIP, metrics.RRoute, metrics.RUserAgent)
	return nil
}

func startMetrics() error {
	if config.Prometheus.Enable {
		return http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), promhttp.Handler())
	}
	return nil
}

func getVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}

func readToml(file string, out tomlFile) error {
	// Preprocessor
	r, err := preproc.Preprocess(file)
	if err != nil {
		return err
	}

	if _, err := toml.DecodeReader(r, out); err != nil {
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
