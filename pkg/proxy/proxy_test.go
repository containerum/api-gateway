package proxy

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStripPath(t *testing.T) {
	Convey("Test strip path func", t, func() {
		Convey("No empty path test", func() {
			url, _ := url.Parse("http://kube-api:1214/namespaces/ivan")
			listenPath := "/kube/namespaces"
			stripPath := stripPath(listenPath, url.Path)
			So(stripPath, ShouldEqual, "/namespaces/ivan")
			So(url.Scheme+"://"+url.Host+stripPath, ShouldEqual, "http://kube-api:1214/namespaces/ivan")
		})
	})
}

func TestBuildHostUrl(t *testing.T) {
	Convey("Test build host url func", t, func() {
		Convey("Good url test", func() {
			u, _ := url.Parse("http://kube-api:1214/namespaces/ivan")
			So(buildHostUrl(*u), ShouldEqual, "http://kube-api:1214")
		})
	})
}
