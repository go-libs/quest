package quest

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrintln(t *testing.T) {
	q, _ := Request(GET, "http://httpbin.org/get")
	Convey("Println request", t, func() {
		So(q.Println(), ShouldEqual, "GET http://httpbin.org/get")
	})
	q, _ = Request(GET, "http://httpbin.org/get")
	q.Do()
	Convey("Println request", t, func() {
		So(q.Println(), ShouldEqual, "GET http://httpbin.org/get "+strconv.Itoa(q.res.StatusCode))
	})
}
