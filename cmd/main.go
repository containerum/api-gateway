package main

import (
	"fmt"
	"os"

	"bitbucket.org/exonch/ch-gateway/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ch-gateway"
	app.Usage = "Awesome Golang API Gateway"
	app.Version = version.String()
	app.Flags = flags
	app.Action = server

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}