package main

import "gopkg.in/urfave/cli.v2"

var (
	DebugFlag = cli.BoolFlag{
		EnvVars: []string{"GATEWAY_DEBUG"},
		Name:    "debug",
		Usage:   "start the server in debug mode",
		Aliases: []string{"d"},
	}

	AuthAddrFlag = cli.StringFlag{
		EnvVars: []string{"GRPC_AUTH_ADDRESS"},
		Name:    "grpc-auth-address",
		Usage:   "GRPC Auth service address",
		Value:   "192.168.88.200",
	}

	ConfigPathFlag = cli.StringFlag{
		EnvVars: []string{"CONFIG_FILE"},
		Name:    "config",
		Usage:   "Path for config TOML file",
		Value:   "config.toml",
	}

	RoutesPathFlag = cli.StringFlag{
		EnvVars: []string{"ROUTES_FILE"},
		Name:    "routes",
		Usage:   "Path for routes TOML file",
		Value:   "routes.toml",
	}

	TLSCertPathFlag = cli.StringFlag{
		EnvVars: []string{"TLS_CERT"},
		Name:    "tls-cert",
		Usage:   "Cert.pem for HTTPS",
		Value:   "cert.pem",
	}

	TLSKeyPathFlag = cli.StringFlag{
		EnvVars: []string{"TLS_KEY"},
		Name:    "tls-key",
		Usage:   "Key.pem for HTTPS",
		Value:   "key.pem",
	}

	ServiceHostPrefixFlag = cli.StringFlag{
		EnvVars: []string{"SERVICE_HOST_PREFIX"},
		Name:    "service-host-prefix",
		Usage:   "Prefix for service hostname (needed for helm)",
		Value:   "",
	}
)
