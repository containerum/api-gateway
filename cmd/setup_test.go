package main

import (
	"testing"

	"git.containerum.net/ch/api-gateway/pkg/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersion(t *testing.T) {
	Convey("Test Get Version function", t, func() {
		Convey("Test default", func() {
			v := getVersion()
			So(v, ShouldNotBeBlank)
			So(v, ShouldEqual, "1.0.0-dev")
		})
		Convey("Test changed", func() {
			Version = "change"
			v := getVersion()
			So(v, ShouldNotBeBlank)
			So(v, ShouldEqual, "change")
		})
	})
}

func TestReadToml(t *testing.T) {
	Convey("Test Read Toml function", t, func() {
		Convey("Empty path", func() {
			var config model.Config
			err := readToml("", &config)
			So(err, ShouldNotBeNil)
		})
		Convey("Wrong path or file", func() {
			var config model.Config
			err := readToml("cli.go", &config)
			So(err, ShouldNotBeNil)
		})
		Convey("Read good config file", func() {
			var config model.Config
			err := readToml("../pkg2/model/tests/good_config.toml", &config)
			So(err, ShouldBeNil)
		})
		Convey("Read good route file", func() {
			var routes model.Routes
			err := readToml("../pkg2/model/tests/good_route.toml", &routes)
			So(err, ShouldBeNil)
		})
		Convey("Read bad config file", func() {
			var config model.Config
			err := readToml("../pkg2/model/tests/bad_config.toml", &config)
			So(err, ShouldNotBeNil)
		})
		Convey("Read bad route file", func() {
			var routes model.Routes
			err := readToml("../pkg2/model/tests/bad_route.toml", &routes)
			So(err, ShouldNotBeNil)
		})
		Convey("Read invalid config file", func() {
			var config model.Config
			err := readToml("../pkg2/model/tests/invalid_config.toml", &config)
			So(err, ShouldNotBeNil)
		})
		Convey("Read invalid route file", func() {
			var routes model.Routes
			err := readToml("../pkg2/model/tests/invalid_route.toml", &routes)
			So(err, ShouldNotBeNil)
		})
	})
}
