package quest

import (
	"testing"

	. "github.com/go-libs/methods"

	mocha "github.com/smartystreets/goconvey/convey"
)

func TestMakeARequest(t *testing.T) {
	q := Request(GET, "http://httpbin.org/get")

	mocha.Convey("Should be making a Request", t, func() {
		mocha.So(q.Method, mocha.ShouldEqual, GET)
	})
}
