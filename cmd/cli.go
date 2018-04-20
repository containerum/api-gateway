package main

import "github.com/urfave/cli"

const usageText = "Awesome Golang API Gateway."

const (
	authAddr    = "grpc-auth-address"
	authPort    = "grpc-auth-port"
	tlsCertPath = "tls-cert"
	tlsKeyPath  = "tls-key"
	routesPath  = "routes"
	configPath  = "config"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GATEWAY_DEBUG",
		Name:   "debug, d",
		Usage:  "start the server in debug mode",
	},
	cli.StringFlag{
		EnvVar: "GRPC_AUTH_ADDRESS",
		Name:   authAddr,
		Usage:  "GRPC Auth service address",
		Value:  "192.168.88.200",
	},
	cli.StringFlag{
		EnvVar: "GRPC_AUTH_PORT",
		Name:   authPort,
		Usage:  "GRPC Auth service port",
		Value:  "1112",
	},
	cli.StringFlag{
		EnvVar: "CONFIG_FILE",
		Name:   configPath,
		Usage:  "Path for config TOML file",
		Value:  "config.toml",
	},
	cli.StringFlag{
		EnvVar: "ROUTES_FILE",
		Name:   routesPath,
		Usage:  "Path for routes TOML file",
		Value:  "routes.toml",
	},
	cli.StringFlag{
		EnvVar: "TLS_CERT",
		Name:   tlsCertPath,
		Usage:  "Cert.pem for HTTPS",
		Value:  "cert.pem",
	},
	cli.StringFlag{
		EnvVar: "TLS_KEY",
		Name:   tlsKeyPath,
		Usage:  "Key.pem for HTTPS",
		Value:  "key.pem",
	},
}
