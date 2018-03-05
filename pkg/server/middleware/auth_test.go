package middleware

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckUserRole(t *testing.T) {
	Convey("Test check user role", t, func() {
		Convey("Without token", func() {
			ok := checkUserRole("admin", []string{})
			So(ok, ShouldBeTrue)
		})
		Convey("For everyone (*)", func() {
			ok := checkUserRole("admin", []string{"user", "*"})
			So(ok, ShouldBeTrue)
			ok = checkUserRole("user", []string{"*"})
			So(ok, ShouldBeTrue)
		})
		Convey("Only user", func() {
			ok := checkUserRole("admin", []string{"user"})
			So(ok, ShouldBeFalse)
			ok = checkUserRole("user", []string{"user"})
			So(ok, ShouldBeTrue)
		})
		Convey("Only admin", func() {
			ok := checkUserRole("admin", []string{"admin"})
			So(ok, ShouldBeTrue)
			ok = checkUserRole("user", []string{"admin"})
			So(ok, ShouldBeFalse)
		})
	})
}
