package model

import (
	"encoding/json"
	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
)

var (
	rateTypes = []string{"local"}
	//ErrInvalidRateType returns when Rate type is not correct
	ErrInvalidRateType = fmt.Errorf("Invalid rate type. Options: %v", rateTypes)
)

type Enable bool

type Config struct {
	Port int
	TLS  struct {
		Enable
		Key  string `toml:"-"`
		Cert string `toml:"-"`
	}
	Auth struct {
		Enable
	}
	Prometheus struct {
		Enable
		Port int
	}
	Rate struct {
		Enable
		Limit int
		Type  string
	}
	Clickhouse struct {
		Enable
		Crypt bool
	}
}

//Validate return array or invalid inputs
//TODO: Add rate limit validation
func (c *Config) Validate() []error {
	var errs []error
	if !valid.IsIn(c.Rate.Type, rateTypes...) {
		errs = append(errs, ErrInvalidRateType)
	}
	return errs
}

func (c *Config) String() string {
	var str []byte
	switch logrus.StandardLogger().Formatter.(type) {
	case *logrus.TextFormatter:
		str, _ = json.MarshalIndent(c, "", "  ")
	default:
		str, _ = json.Marshal(c)
	}
	return string(str)
}
