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
	app.Version = GetVersion()
	app.Flags = flags
	app.Usage = usageText
	app.Action = runServer
	app.Before = setLogFormat

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
