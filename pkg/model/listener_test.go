package model

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	l1 = Listener{
		DefaultModel: DefaultModel{
			ID:        "11472afa-48ff-4583-b2e4-74119993f22a",
			CreatedAt: time.Date(2017, 12, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 12, 1, 12, 40, 0, 0, time.Local),
		},
		Name:        "Router N1",
		OAuth:       newBool(false),
		Active:      newBool(false),
		Group:       g1,
		GroupRefer:  "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
		StripPath:   newBool(false),
		ListenPath:  "route1",
		UpstreamURL: "http://localhost",
		Method:      "GET",
	}
	g1 = Group{
		DefaultModel: DefaultModel{
			ID:        "7c7ebd79-8e9d-443e-8c20-d43de4e33b45",
			CreatedAt: time.Date(2017, 13, 1, 12, 30, 0, 0, time.Local),
			UpdatedAt: time.Date(2017, 16, 1, 13, 50, 0, 0, time.Local),
		},
		Name: "Group N1",
	}
)

var (
	listenerJSON = []string{
		//0. Good server answer. Full marshal fields
		`{"id":"11472afa-48ff-4583-b2e4-74119993f22a","created_at":1512120600,` +
			`"updated_at":1512121200,"name":"Router N1","o_auth":false,"active":false,` +
			`"group":{"id":"7c7ebd79-8e9d-443e-8c20-d43de4e33b45","created_at":1514799000,` +
			`"updated_at":1522579800,"name":"Group N1"},"strip_path":false,` +
			`"listen_path":"route1","upstream_url":"http://localhost","method":"GET"}`,
		//1. Good client request for create. Full fields.
		`{"name":"Router N1","o_auth":false,"active":false,` +
			`"group_id":"7c7ebd79-8e9d-443e-8c20-d43de4e33b45","strip_path":false,` +
			`"listen_path":"route1","upstream_url":"http://localhost","method":"GET"}`,
		//2. Good client request for update. Full fields.
		`{"name":"Router N1","o_auth":false,"active":false,` +
			`"group_id":"7c7ebd79-8e9d-443e-8c20-d43de4e33b45","strip_path":false,` +
			`"listen_path":"route1","upstream_url":"http://localhost","method":"GET"}`,
		//3. Good client request for update active. Full fields.
		`{"active":true}`,
		//4. Bad client request. Empty fields.
		`{}`,
		//5. Bad client request. Wrong fields.
		`{"wrong":"data"}`,
		//6. Bad client request for update listener. Wrong active type
		`{"active":"baddata"}`,
		//7. Bad client request for update listener. Wrong oauth.
		`{"o_auth":"wrong"}`,
		//8. Bad client request for update and create listener. Invalid name length
		`{"name":"x"}`,
		//9. Bad client request for update and create listener. Invalid name length
		`{"name":"Y8tkEM01S7E69NJ4onpJs1y0LPDcWmqkyCR6BKvhyWm3lVQHFehNbaH1kdTMxv0ocl2vPz4bXmB2g5qc2TuKH6xoHTt8fdiYtl1cezxBa3CPXGYZt2nfsJzv1gHeEWarKv"}`,
		//10. Bad client request for update and create listener. Invalid method
		`{"name":"Good name", "method":"BADMETHOD"}`,
		//11. Bad client request for update and create listener. Invalid group id
		`{"name":"Good name", "group_id":"11472afa-48ff-4583"}`,
		//12. Bad client request for update and create listener. Invalid listen path length
		`{"name":"Good name", "listen_path":"y"}`,
		//13. Bad client request for update and create listener. Invalid listen path length
		`{"name":"Good name", "listen_path":"Y8tkEM01S7E69NJ4onpJs1y0LPDcWmqkyCR6BKvhyWm3lVQHFehNbaH1kdTMxv0ocl2vPz4bXmB2g5qc2TuKH6xoHTt8fdiYtl1cezxBa3CPXGYZt2nfsJzv1gHeEWarKv"}`,
		//14. Bad client request for update and create listener. Invalid upstream url length
		`{"name":"Good name", "upstream_url":"y"}`,
		//15. Bad client request for update and create listener. Invalid upstream url length
		`{"name":"Good name", "upstream_url":"Y8tkEM01S7E69NJ4onpJs1y0LPDcWmqkyCR6BKvhyWm3lVQHFehNbaH1kdTMxv0ocl2vPz4bXmB2g5qc2TuKH6xoHTt8fdiYtl1cezxBa3CPXGYZt2nfsJzv1gHeEWarKv"}`,
		//16. Bad client request for create. Empty group id
		`{"name":"Router N1","o_auth":false,"active":false,"strip_path":false,` +
			`"listen_path":"route1","upstream_url":"http://localhost","method":"GET"}`,
		//17. Good client request for update oauth. Full fields.
		`{"o_auth":true}`,
	}
)

func TestMarshalJSON(t *testing.T) {
	Convey("Listener JSON Marshal test", t, func() {
		Convey("Test marshal full field listener", func() {
			b, err := l1.MarshalJSON()
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, listenerJSON[0])
		})
	})
}

func TestUnmarshalJSON(t *testing.T) {
	Convey("Listener JSON Unmarshal test", t, func() {
		Convey("Wrong data unmarshal", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[5]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.Name, ShouldBeBlank)
			So(listener.GroupRefer, ShouldBeBlank)
			So(listener.ListenPath, ShouldBeBlank)
			So(listener.UpstreamURL, ShouldBeBlank)
			So(listener.Method, ShouldBeBlank)
		})
		Convey("Empty json Unmarshal", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[4]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.Name, ShouldBeBlank)
			So(listener.GroupRefer, ShouldBeBlank)
			So(listener.ListenPath, ShouldBeBlank)
			So(listener.UpstreamURL, ShouldBeBlank)
			So(listener.Method, ShouldBeBlank)
		})
		Convey("Test Unmashal for Create Listener", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[1]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.Name, ShouldNotBeBlank)
			So(listener.GroupRefer, ShouldNotBeBlank)
			So(listener.ListenPath, ShouldNotBeBlank)
			So(listener.UpstreamURL, ShouldNotBeBlank)
			So(listener.Method, ShouldNotBeBlank)
		})
		Convey("Test Unmarshal for Update", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[2]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.Name, ShouldNotBeBlank)
			So(listener.GroupRefer, ShouldNotBeBlank)
			So(listener.ListenPath, ShouldNotBeBlank)
			So(listener.UpstreamURL, ShouldNotBeBlank)
			So(listener.Method, ShouldNotBeBlank)
		})
		Convey("Test Unmarshal for Update Activate", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[3]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.Active, ShouldNotBeNil)
		})
		Convey("Test Unmarshal for Update OAuth", func() {
			listener := Listener{}
			err := listener.UnmarshalJSON([]byte(listenerJSON[17]))
			So(err, ShouldBeNil)
			So(listener.ID, ShouldBeBlank)
			So(listener.OAuth, ShouldNotBeNil)
		})
	})
}

func TestValidateCreate(t *testing.T) {
	Convey("Listener Validate Create test", t, func() {
		Convey("Good request", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[1]))
			err := listener.ValidateCreate()
			So(err, ShouldBeEmpty)
		})
		Convey("Bad request", func() {
			Convey("Check min name length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[8]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 5)
				So(err, ShouldContain, ErrInvalidListenerNameLength)
			})
			Convey("Check max name length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[9]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 5)
				So(err, ShouldContain, ErrInvalidListenerNameLength)
			})
			Convey("Check bad method", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[10]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerMethod)
			})
			Convey("Check bad group id", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[11]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerGroupRefer)
			})
			Convey("Check min listen path length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[12]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerListenPath)
			})
			Convey("Check max listen path length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[13]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerListenPath)
			})
			Convey("Check min upstream url length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[14]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerUpstreamURL)
			})
			Convey("Check max upstream url length", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[15]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 4)
				So(err, ShouldContain, ErrInvalidListenerUpstreamURL)
			})
			Convey("Check empty group id", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[16]))
				err := listener.ValidateCreate()
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 1)
				So(err, ShouldContain, ErrInvalidListenerGroupRefer)
			})
		})
	})
}

func TestValidateUpdate(t *testing.T) {
	Convey("Listener Validate Update test", t, func() {
		Convey("Check create model", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[1]))
			err := listener.ValidateUpdate("")
			So(err, ShouldNotBeEmpty)
			So(err, ShouldContain, ErrInvalidListenerID)
		})
		Convey("Bad request", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[6]))
			err := listener.ValidateUpdate("")
			So(err, ShouldNotBeEmpty)
			So(err, ShouldHaveLength, 6)
			So(err, ShouldContain, ErrInvalidListenerID)
			So(err, ShouldContain, ErrInvalidListenerMethod)
		})
		Convey("Good request", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[2]))
			err := listener.ValidateUpdate("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(err, ShouldBeEmpty)
		})
	})
}

func TestValidateUpdateActive(t *testing.T) {
	Convey("Listener Validate Update Active test", t, func() {
		Convey("Nil listener", func() {
			listener := &Listener{}
			listener = nil
			err := listener.ValidateUpdateActive("")
			So(err, ShouldNotBeEmpty)
			So(err, ShouldHaveLength, 1)
			So(err, ShouldContain, ErrNilListener)
		})
		Convey("Good request", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[3]))
			err := listener.ValidateUpdateActive("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(err, ShouldBeEmpty)
		})
		Convey("Bad requests", func() {
			Convey("Empty request", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[4]))
				err := listener.ValidateUpdateActive("")
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 2)
				So(err, ShouldContain, ErrInvalidListenerID)
				So(err, ShouldContain, ErrInvalidListenerActive)
			})
			Convey("Empty ID", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[6]))
				err := listener.ValidateUpdateActive("")
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 2)
				So(err, ShouldContain, ErrInvalidListenerID)
			})
			Convey("Empty Active", func() {
				listener := Listener{}
				listener.UnmarshalJSON([]byte(listenerJSON[7]))
				err := listener.ValidateUpdateActive("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
				So(err, ShouldNotBeEmpty)
				So(err, ShouldHaveLength, 1)
				So(err, ShouldContain, ErrInvalidListenerActive)
			})
		})
	})
}

func TestGetUpdateType(t *testing.T) {
	Convey("Test update method chose", t, func() {
		Convey("Update None", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[14]))
			update, err := listener.GetUpdateType("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(update, ShouldEqual, ListenerUpdateNone)
			So(err, ShouldNotBeEmpty)
		})
		Convey("Update Full", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[2]))
			update, err := listener.GetUpdateType("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(update, ShouldEqual, ListenerUpdateFull)
			So(err, ShouldBeEmpty)
		})
		Convey("Update Active", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[3]))
			update, err := listener.GetUpdateType("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(update, ShouldEqual, ListenerUpdateActive)
			So(err, ShouldBeEmpty)
		})
		Convey("Update OAuth", func() {
			listener := Listener{}
			listener.UnmarshalJSON([]byte(listenerJSON[17]))
			update, err := listener.GetUpdateType("7c7ebd79-8e9d-443e-8c20-d43de4e33b45")
			So(update, ShouldEqual, ListenerUpdateOAuth)
			So(err, ShouldBeEmpty)
		})
	})
}

//TODO: Move to utils
func newBool(b bool) *bool {
	res := new(bool)
	*res = b
	return res
}
