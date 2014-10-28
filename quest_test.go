package quest

import (
	"fmt"
	"net/http"
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

func TestResponseHandling(t *testing.T) {
	Request(GET, "http://httpbin.org/get").
		Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	})
}
