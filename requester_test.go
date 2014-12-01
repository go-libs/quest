package quest

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func TestDebugPrintln(t *testing.T) {
	c1 := &http.Cookie{}
	c1.Name = "k1"
	c1.Value = "v1"
	c2 := &http.Cookie{}
	c2.Name = "k2"
	c2.Value = "v2"
	queryParams := url.Values{}
	queryParams.Set("foo", "bar")
	queryParams.Set("name", "bazz")
	q, _ := Request(GET, "http://httpbin.org/cookies")
	q.Query(&queryParams)
	q.Cookie(c1, c2)
	Convey("DebugPrintln request", t, func() {
		s := []string{"$ curl -i", "-b " + strconv.Quote("k1=v1; k2=v2"), strconv.Quote("http://httpbin.org/cookies?foo=bar&name=bazz")}
		So(q.DebugPrintln(), ShouldEqual, strings.Join(s, " \\\n\t"))
	})
}
