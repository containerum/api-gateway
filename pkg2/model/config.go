package model

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
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
