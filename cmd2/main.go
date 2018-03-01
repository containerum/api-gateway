package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

//Version keeps curent app version
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

//TODO: Rewrite IF noodles
func runServer(c *cli.Context) error {
	setupLogs(c)
	if err := setupConfig(c); err != nil {
		return err
	}
	if err := setupRoutes(c); err != nil {
		return err
	}
	if err := setupTLS(c); err != nil {
		return err
	}
	if err := setupAuth(c); err != nil {
		return err
	}
	serv, err := setupServer(c)
	if err != nil {
		return err
	}
	return serv.Start()
}
