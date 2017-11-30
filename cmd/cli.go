package main

import "github.com/urfave/cli"

const usageText = "Awesome Golang API Gateway. \"Not Safety\" migrations runs only migrate mode!"

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GATEWAY_DEBUG",
		Name:   "debug, d",
		Usage:  "start the server in debug mode",
	},
	cli.BoolFlag{
		Name:  "migrate, m",
		Usage: "start the server in migration mode",
	},
	cli.StringFlag{
		EnvVar: "PG_USER",
		Name:   "pg-user",
		Usage:  "Postgres user",
		Value:  "pg",
	},
	cli.StringFlag{
		EnvVar: "PG_PASSWORD",
		Name:   "pg-password",
		Usage:  "Postgres user password",
		Value:  "123456789",
	},
	cli.StringFlag{
		EnvVar: "PG_DATABASE",
		Name:   "pg-database",
		Usage:  "Postgres database",
		Value:  "postgres",
	},
	cli.StringFlag{
		EnvVar: "PG_ADDRESS",
		Name:   "pg-address",
		Usage:  "Postgres address",
		Value:  "x1.containerum.io",
	},
	cli.StringFlag{
		EnvVar: "PG_PORT",
		Name:   "pg-port",
		Usage:  "Postgres port",
		Value:  "36519",
	},
	cli.BoolFlag{
		EnvVar: "PG_DEBUG",
		Name:   "pg-debug",
		Usage:  "Write gorm logs",
	},
	cli.BoolFlag{
		EnvVar: "PG_SAFE_MIGRATION",
		Name:   "pg-safe-migration",
		Usage:  "Use safe migration",
	},
	cli.StringFlag{
		EnvVar: "STATSD_ADDRESS",
		Name:   "statsd-address",
		Usage:  "Statsd address",
		Value:  "213.239.208.25:8125",
	},
	cli.StringFlag{
		EnvVar: "STATSD-PREFIX",
		Name:   "statsd-prefix",
		Usage:  "Statsd prefix",
		Value:  "ch-gateway",
	},
	cli.IntFlag{
		EnvVar: "STATSD-BUFFER-TIME",
		Name:   "statsd-buffer-time",
		Usage:  "Statsd buffer time",
		Value:  300,
	},
}

var commands = []cli.Command{
	cli.Command{
		Name:        "migration",
		Subcommands: migrationCommands,
	},
}

var migrationCommands = []cli.Command{
	cli.Command{
		Name:   "init",
		Usage:  "Creates migrations table",
		Flags:  flags,
		Before: setLogFormat,
		Action: initMigration,
	},
	cli.Command{
		Name:   "version",
		Usage:  "Prints current db version",
		Flags:  flags,
		Before: setLogFormat,
		Action: getMigrationVersion,
	},
	cli.Command{
		Name:   "up",
		Usage:  "Runs all available migrations",
		Flags:  flags,
		Before: setLogFormat,
		Action: upMigration,
	},
	cli.Command{
		Name:   "down",
		Usage:  "Reverts last migration",
		Flags:  flags,
		Before: setLogFormat,
		Action: downMigration,
	},
}
