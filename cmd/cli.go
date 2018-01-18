package main

import "github.com/urfave/cli"

const usageText = "Awesome Golang API Gateway. \"Not Safety\" migrations runs only migrate mode!"

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GATEWAY_DEBUG",
		Name:   "debug, d",
		Usage:  "start the server in debug mode",
	},
	cli.StringFlag{
		EnvVar: "PG_USER",
		Name:   "pg-user",
		Usage:  "Postgres user",
		Value:  "gatewayapi",
	},
	cli.StringFlag{
		EnvVar: "PG_PASSWORD",
		Name:   "pg-password",
		Usage:  "Postgres user password",
		Value:  "Caik9miegh6k",
	},
	cli.StringFlag{
		EnvVar: "PG_DATABASE",
		Name:   "pg-database",
		Usage:  "Postgres database",
		Value:  "gatewayapi",
	},
	cli.StringFlag{
		EnvVar: "PG_ADDRESS",
		Name:   "pg-address",
		Usage:  "Postgres address",
		Value:  "192.168.88.200",
	},
	cli.StringFlag{
		EnvVar: "PG_PORT",
		Name:   "pg-port",
		Usage:  "Postgres port",
		Value:  "5432",
	},
	cli.BoolFlag{
		EnvVar: "PG_MIGRATIONS",
		Name:   "pg-migrations",
		Usage:  "Run migrations",
	},
	cli.StringFlag{
		EnvVar: "STATSD_ADDRESS",
		Name:   "statsd-address",
		Usage:  "Statsd address",
		Value:  "192.168.88.200:8125",
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
	cli.StringFlag{
		EnvVar: "GRPC_AUTH_ADDRESS",
		Name:   "grpc-auth-address",
		Usage:  "GRPC Auth service address",
		Value:  "192.168.88.200",
	},
	cli.StringFlag{
		EnvVar: "GRPC_AUTH_PORT",
		Name:   "grpc-auth-port",
		Usage:  "GRPC Auth service port",
		Value:  "1112",
	},
	cli.StringFlag{
		EnvVar: "REDIS_ADDRESS",
		Name:   "redis-address",
		Usage:  "Redis service address",
		Value:  "192.168.88.200:6379",
	},
	cli.StringFlag{
		EnvVar: "REDIS_PASSWORD",
		Name:   "redis-password",
		Usage:  "Redis service password",
		Value:  "",
	},
	cli.StringFlag{
		EnvVar: "RATE_LIMIT",
		Name:   "rate-limit",
		Usage:  "Limit requests per second",
		Value:  "3",
	},
	cli.StringFlag{
		EnvVar: "CLICKHOUSE_LOGGER",
		Name:   "clickhouse-logger",
		Usage:  "Write all logs to clickhouse",
		Value:  "88.99.160.131:7777",
	},
	cli.StringFlag{
		EnvVar: "TLS_CERT",
		Name:   "tls-cert",
		Usage:  "Cert.pem for HTTPS",
		Value:  "cert.pem",
	},
	cli.StringFlag{
		EnvVar: "TLS_KEY",
		Name:   "tls-key",
		Usage:  "Key.pem for HTTPS",
		Value:  "key.pem",
	},
}
