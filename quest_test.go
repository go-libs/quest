package quest

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"os"
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

	mocha.Convey("QueryParameters, query string", t, func() {
		Request(GET, "http://httpbin.org/get").
			QueryParameters(&queryParams).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			mocha.So(req.URL.String(), mocha.ShouldEqual, "http://httpbin.org/get?foo=bar&name=%E6%B4%BB%E5%8A%9B")
		})
	})

	mocha.Convey("Parameters, ContentLength should equal to buffer length", t, func() {
		Request(POST, "http://httpbin.org/post").
			Parameters(queryParams).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			mocha.So(res.ContentLength, mocha.ShouldEqual, int64(data.Len()))
		})
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

	mocha.Convey("Response JSON", t, func() {
		Request(POST, "http://httpbin.org/post").
			Encoding("JSON").
			Parameters(parameters).
			ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
			mocha.Convey("Data - a pointer struct", func() {
				mocha.So(data, mocha.ShouldPointTo, data)
				mocha.So(data.Headers["Host"], mocha.ShouldEqual, "httpbin.org")
			})
		}).
			ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct2, e error) {
			mocha.Convey("Data - a struct", func() {
				mocha.So(&data, mocha.ShouldNotPointTo, &DataStruct2{})
				mocha.So(data.Origin, mocha.ShouldNotBeNil)
			})
		}).
			// Nothing happend
			ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct2, e error, g error) {
			log.Println("Nothing happend!")
		})
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

	mocha.Convey("Response JSON, using JSON decode", t, func() {
		Request(POST, "http://httpbin.org/post").
			Encoding("JSON").
			Parameters(parameters2).
			ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct4, e error) {
			mocha.Convey("Using DataStruct4 JSON struct", func() {
				mocha.So(data.Origin, mocha.ShouldNotBeNil)
			})
		}).
			ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct3, e error) {
			mocha.Convey("Using DataStruct3 JSON struct", func() {
				mocha.So(&data.Json, mocha.ShouldResemble, parameters2)
			})
		})
	})

	mocha.Convey("Encoding Query Options", t, func() {
		type Options struct {
			Foo string `url:"foo"`
			Baz []int  `url:"baz"`
		}

		// http://httpbin.org/get
		Request(GET, "http://httpbin.org/get").
			QueryParameters(Options{"bar", []int{233, 377, 610}}).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			mocha.So(req.URL.String(), mocha.ShouldEqual, "http://httpbin.org/get?baz=233&baz=377&baz=610&foo=bar")
		})
	})
}

func TestDownload(t *testing.T) {
	os.Mkdir("tmp", os.ModePerm)

	mocha.Convey("Downloading file", t, func() {
		mocha.Convey("Downloading stream.log in progress", func() {
			Download(GET, "http://httpbin.org/bytes/1024", "tmp/stream.log").
				Progress(func(c, t, e int64) {
				log.Println(c, t, e)
				mocha.So(c, mocha.ShouldBeLessThanOrEqualTo, t)
			}).Do()
		})
		mocha.Convey("Downloading stream2.log in progress and invoke response handler", func() {
			var n int64
			Download(GET, "http://httpbin.org/bytes/10240", "tmp/stream2.log").
				Progress(func(c, t, e int64) {
				n = c
				log.Println(c, t, e)
				mocha.So(c, mocha.ShouldBeLessThanOrEqualTo, t)
			}).Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
				l := int64(data.Len())
				mocha.So(n, mocha.ShouldEqual, l)
				mocha.So(res.ContentLength, mocha.ShouldEqual, l)
			})
		})
	})
}

func TestUpload(t *testing.T) {
	mocha.Convey("Uploading one file", t, func() {
		Upload(POST, "http://httpbin.org/post", map[string]string{"stream": "tmp/stream.log"}, nil).
			Progress(func(c, tt, e int64) {
			log.Println(c, tt, e)
			//mocha.So(c, mocha.ShouldBeLessThanOrEqualTo, tt)
		}).Do()
	})
	mocha.Convey("Uploading multi files", t, func() {
		Upload(POST, "http://httpbin.org/post", map[string]string{"stream": "tmp/stream.log", "stream2": "tmp/stream2.log"}, nil).
			Progress(func(c, tt, e int64) {
			log.Println(c, tt, e)
			//mocha.So(c, mocha.ShouldBeLessThanOrEqualTo, tt)
		}).Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			l := int64(data.Len())
			mocha.So(res.ContentLength, mocha.ShouldEqual, l)
		})
	})
}
