package model

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
)

var (
	roles         = []string{"*", "user", "admin"}
	methods       = []string{"GET", "POST", "PUT", "DELETE"}
	upstreamProto = "http"
	//ErrInvalidRoles returns when roles is wrong
	ErrInvalidRoles = fmt.Errorf("Invalid roles. Options: %v", roles)
	//ErrInvalidMethod returns when method is wrong
	ErrInvalidMethod = fmt.Errorf("Invalid method. Options: %v", methods)
	//ErrInvalidUpstreamProtocol return when usptream is incorrect
	ErrInvalidUpstreamProtocol = fmt.Errorf("Invalid usptream protocol. Options: %v", upstreamProto)
	//ErrInvalidListen returns when listen path is wrong
	ErrInvalidListen = fmt.Errorf("Invalid listen path.")
)

type Routes struct {
	Routes map[string]Route
}

type Route struct {
	ID       string `toml:"-"`
	Name     string
	Desc     string
	Active   bool
	Roles    []string
	Method   string
	Upstream string
	Listen   string
	Strip    bool
	Group    string
	WS       bool
}

//Validate return array or invalid inputs
func (rs *Routes) Validate() []error {
	var errs []error
	for _, route := range rs.Routes {
		//Validate Roles
		for _, role := range route.Roles {
			if !valid.IsIn(role, roles...) {
				errs = append(errs, ErrInvalidRateType)
			}
		}
		//Validate Method
		if !valid.IsIn(route.Method, methods...) {
			errs = append(errs, ErrInvalidMethod)
		}
		// //Validate Upstream protocol
		// if !valid.Matches(route.Upstream, "^"+upstreamProto+"://") {
		// 	errs = append(errs, ErrInvalidUpstreamProtocol)
		// }
		//Validate Listen
		if !valid.Matches(route.Listen, "^/") {
			errs = append(errs, ErrInvalidListen)
		}
	}
	return errs
}
