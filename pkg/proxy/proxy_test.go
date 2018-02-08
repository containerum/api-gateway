package proxy

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStripPath(t *testing.T) {
	Convey("Test strip path func", t, func() {
		Convey("Listen path and request path equals and not empty", func() {
			req, _ := url.Parse("http://localhost:1313/kube/namespaces")
			upstream, _ := url.Parse("http://kube-api:1214/namespaces")
			listenPath := "/kube/namespaces"
			stripPath := stripPath(req.Path, listenPath, upstream.Path)
			So(stripPath, ShouldEqual, "/namespaces")
		})
		Convey("Listen path and request path not equals and not empty", func() {
			req, _ := url.Parse("http://localhost:1313/kube/namespaces/ns1")
			upstream, _ := url.Parse("http://kube-api:1214/namespaces")
			listenPath := "/kube/namespaces/*"
			stripPath := stripPath(req.Path, listenPath, upstream.Path)
			So(stripPath, ShouldEqual, "/namespaces/ns1")
		})
		Convey("Large path", func() {
			req, _ := url.Parse("http://localhost:1313/kube/namespaces/ns1/deploy/x2")
			upstream, _ := url.Parse("http://kube-api:1214/namespaces")
			listenPath := "/kube/namespaces/*"
			stripPath := stripPath(req.Path, listenPath, upstream.Path)
			So(stripPath, ShouldEqual, "/namespaces/ns1/deploy/x2")
		})
		Convey("Empty upstream", func() {
			req, _ := url.Parse("http://localhost:1313/kube/namespaces/ns1/deploy/x2")
			upstream, _ := url.Parse("http://kube-api:1214")
			listenPath := "/kube/namespaces/*"
			stripPath := stripPath(req.Path, listenPath, upstream.Path)
			So(stripPath, ShouldEqual, "/ns1/deploy/x2")
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
