package main

import "gopkg.in/urfave/cli.v2"

const usageText = "Awesome Golang API Gateway."

const (
	authAddr          = "grpc-auth-address"
	tlsCertPath       = "tls-cert"
	tlsKeyPath        = "tls-key"
	routesPath        = "routes"
	configPath        = "config"
	serviceHostPrefix = "service-host-prefix"
)

var (
	DebugFlag = cli.BoolFlag{
		EnvVars: []string{"GATEWAY_DEBUG"},
		Name:    "debug",
		Usage:   "start the server in debug mode",
		Aliases: []string{"d"},
	}

	AuthAddrFlag = cli.StringFlag{
		EnvVars: []string{"GRPC_AUTH_ADDRESS"},
		Name:    authAddr,
		Usage:   "GRPC Auth service address",
		Value:   "192.168.88.200",
	}

	ConfigPathFlag = cli.StringFlag{
		EnvVars: []string{"CONFIG_FILE"},
		Name:    configPath,
		Usage:   "Path for config TOML file",
		Value:   "config.toml",
	}

	RoutesPathFlag = cli.StringFlag{
		EnvVars: []string{"ROUTES_FILE"},
		Name:    routesPath,
		Usage:   "Path for routes TOML file",
		Value:   "routes.toml",
	}

	TLSCertPathFlag = cli.StringFlag{
		EnvVars: []string{"TLS_CERT"},
		Name:    tlsCertPath,
		Usage:   "Cert.pem for HTTPS",
		Value:   "cert.pem",
	}

	TLSKeyPathFlag = cli.StringFlag{
		EnvVars: []string{"TLS_KEY"},
		Name:    tlsKeyPath,
		Usage:   "Key.pem for HTTPS",
		Value:   "key.pem",
	}

	ServiceHostPrefixFlag = cli.StringFlag{
		EnvVars: []string{"SERVICE_HOST_PREFIX"},
		Name:    serviceHostPrefix,
		Usage:   "Prefix for service hostname (needed for helm)",
		Value:   "",
	}
)
