package main

import (
	"fmt"
	"net/http"
	"os"

	"git.containerum.net/ch/api-gateway/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Name = "ch-gateway"
	app.Version = getVersion()
	app.Flags = flags
	app.Usage = usageText
	app.Action = runServer

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServer(c *cli.Context) error {
	var serv *server.Server
	var err error
	if err = setup(c, setupLogs, setupConfig, setupRoutes, setupTLS, setupAuth, setupMetrics); err != nil {
		return err
	}
	if serv, err = setupServer(c); err != nil {
		return err
	}
	go startMetrics()
	return serv.Start()
}

func getVersion() string {
	if Version == "" {
		return "1.0.0-dev"
	}
	return Version
}

func startMetrics() error {
	if config.Prometheus.Enable {
		return http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), promhttp.Handler())
	}
	return nil
}
