package quest

import (
	"fmt"
	"net/http"
	"net/url"
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
	queryParams := url.Values{}
	queryParams.Set("foo", "bar")
	queryParams.Set("name", "活力")

	Request(GET, "http://httpbin.org/get").
		Query(queryParams).
		Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	})

	Request(POST, "http://httpbin.org/post").
		Parameters(queryParams).
		Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	})

	parameters := map[string]interface{}{
		"foo": []int{1, 2, 3},
		"bar": map[string]string{"baz": "qux"},
	}

	Request(POST, "http://httpbin.org/post").
		Encoding("JSON").
		Parameters(parameters).
		ResponseJSON(func(request *http.Request, response *http.Response, data JSONMaps, err error) {
		fmt.Println("Response JSON Format")
		fmt.Printf("%+v\n", data["data"])
	})
}
