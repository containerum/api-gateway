package model

type Routes struct {
	Routes map[string]Route
}

type Route struct {
	Name     string
	Desc     string
	Active   bool
	Roles    []string
	Rate     int
	Method   string
	Upstream string
	Listen   string
	Strip    bool
	Group    string
}

func (rs *Routes) Validate() error {

	return nil
}
