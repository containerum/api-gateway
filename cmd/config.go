package main

import (
	"github.com/BurntSushi/toml"
)

//Config strcut for TOML configuration file
type Config struct {
	TLS struct {
		Enable bool
	}
	Store struct {
		Enable bool
	}
	Auth struct {
		Enable bool
	}
	Prometheus struct {
		Enable bool
	}
	Rate struct {
		Enable bool
		Type   string
	}
	Clickhouse struct {
		Enable bool
		Crypt  bool
	}
}

func getConfig() (Config, error) {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		return config, err
	}
	return config, nil
}
