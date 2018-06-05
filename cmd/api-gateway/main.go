package main

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/urfave/cli.v2"
)

var Version string

func main() {
	app := cli.App{
		Name: "api-gateway",
		Version: func() string {
			if Version == "" {
				return "1.0.0-dev"
			}
			return Version
		}(),
		Usage: "Awesome Golang API Gateway.",
		Flags: []cli.Flag{
			&DebugFlag,
			&AuthAddrFlag,
			&ConfigPathFlag,
			&RoutesPathFlag,
			&TLSCertPathFlag,
			&TLSKeyPathFlag,
			&ServiceHostPrefixFlag,
		},
		Before: prettyPrintFlags,
		Action: runServer,
	}

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
	go startMetrics(c)
	return serv.Start()
}

func startMetrics(c *cli.Context) error {
	config := c.App.Metadata[configKey].(*model.Config)
	if config.Prometheus.Enable {
		return http.ListenAndServe(fmt.Sprintf(":%d", config.Prometheus.Port), promhttp.Handler())
	}
	return nil
}

func prettyPrintFlags(ctx *cli.Context) error {
	fmt.Printf("Starting %v %v\n", ctx.App.Name, ctx.App.Version)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)
	for _, f := range ctx.App.VisibleFlags() {
		fmt.Fprintf(w, "Flag: %s\t Value: %v\n", f.Names()[0], ctx.Generic(f.Names()[0]))
	}
	return w.Flush()
}
