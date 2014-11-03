package quest

import (
	"bytes"
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
		QueryParameters(&queryParams).
		Response(func(request *http.Request, response *http.Response, data *bytes.Buffer, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	})

	Request(POST, "http://httpbin.org/post").
		Parameters(queryParams).
		Response(func(request *http.Request, response *http.Response, data *bytes.Buffer, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	})

	parameters := map[string]interface{}{
		"foo": []int{1, 2, 3},
		"bar": map[string]string{"baz": "qux"},
	}

	type DataStruct struct {
		Headers map[string]string
		Origin  string
	}

	type DataStruct2 struct {
		Origin string
	}

	Request(POST, "http://httpbin.org/post").
		Encoding("JSON").
		Parameters(parameters).
		ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
		fmt.Println(data.Headers["Host"])
		fmt.Println(data.Headers["Content-Type"])
		fmt.Println(data.Origin)
	}).
		ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct2, e error) {
		fmt.Println(data.Origin)
	}).
		// Nothing happend
		ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct2, e error, g error) {
		fmt.Println(e)
	})

	type PostParameters struct {
		Foo []int             `json:"foo,omitempty"`
		Bar map[string]string `json:"bar,omitempty"`
	}

	parameters2 := &PostParameters{
		Foo: []int{1, 2, 3},
		Bar: map[string]string{"baz": "qux"},
	}

	type DataStruct4 struct {
		Origin string
	}
	type DataStruct3 struct {
		Headers map[string]string
		Origin  string
		Json    PostParameters `json:"json,omitempty"`
	}

	Request(POST, "http://httpbin.org/post").
		Encoding("JSON").
		Parameters(parameters2).
		ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct4, e error) {
		fmt.Println(data)
	}).
		ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct3, e error) {
		fmt.Println(data)
		fmt.Println(data.Json)
	})

	type Options struct {
		Foo string `url:"foo"`
		Baz []int  `url:"baz"`
	}

	// http://httpbin.org/get http://httpbin.org/get?baz=233&baz=377&baz=610&foo=bar
	fmt.Println(Request(GET, "http://httpbin.org/get").
		QueryParameters(Options{"bar", []int{233, 377, 610}}))
}
