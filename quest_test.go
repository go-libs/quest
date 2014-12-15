package quest

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeARequest(t *testing.T) {
	q, _ := Request(GET, "http://httpbin.org/get")
	Convey("Should be making a Request", t, func() {
		So(q.Method, ShouldEqual, GET)
	})
}

func TestResponseHandling(t *testing.T) {
	queryParams := url.Values{}
	queryParams.Set("foo", "bar")
	queryParams.Set("name", "活力")

	Convey("Query, query string", t, func() {
		q, _ := Request(GET, "http://httpbin.org/get")
		q.
			Query(&queryParams).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			So(req.URL.String(), ShouldEqual, "http://httpbin.org/get?foo=bar&name=%E6%B4%BB%E5%8A%9B")
		})
	})

	Convey("Parameters, ContentLength should equal to buffer length", t, func() {
		q, _ := Request(POST, "http://httpbin.org/post")
		q.
			Parameters(queryParams).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			So(res.ContentLength, ShouldEqual, int64(data.Len()))
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

	Convey("Response JSON", t, func() {
		q, _ := Request(POST, "http://httpbin.org/post")
		q.
			Encoding("JSON").
			Parameters(parameters).
			ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
			Convey("Data - a pointer struct", func() {
				So(data, ShouldPointTo, data)
				So(data.Headers["Host"], ShouldEqual, "httpbin.org")
			})
		}).
			ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct2, e error) {
			Convey("Data - a struct", func() {
				So(&data, ShouldNotPointTo, &DataStruct2{})
				So(data.Origin, ShouldNotBeNil)
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

	Convey("Response JSON, using JSON decode", t, func() {
		q, _ := Request(POST, "http://httpbin.org/post")
		q.
			Encoding("JSON").
			Parameters(parameters2).
			ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct4, e error) {
			Convey("Using DataStruct4 JSON struct", func() {
				So(data.Origin, ShouldNotBeNil)
			})
		}).
			ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct3, e error) {
			Convey("Using DataStruct3 JSON struct", func() {
				So(&data.Json, ShouldResemble, parameters2)
			})
		})
	})

	Convey("Encoding Query Options", t, func() {
		type Options struct {
			Foo string `url:"foo"`
			Baz []int  `url:"baz"`
		}

		// http://httpbin.org/get
		q, _ := Request(GET, "http://httpbin.org/get")
		q.
			Query(Options{"bar", []int{233, 377, 610}}).
			Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
			So(req.URL.String(), ShouldEqual, "http://httpbin.org/get?baz=233&baz=377&baz=610&foo=bar")
		})
	})
}

func TestDownload(t *testing.T) {
	os.Mkdir("tmp", os.ModePerm)

	Convey("Downloading file", t, func() {
		Convey("Downloading stream.log in progress\n", func() {
			q, _ := Download(GET, "http://httpbin.org/bytes/1024", "tmp/stream.log")
			q.
				Progress(func(c, t, e int64) {
				log.Println(c, t, e)
				So(c, ShouldBeLessThanOrEqualTo, t)
			}).Do()
		})
		Convey("Downloading stream2.log in progress and invoke response handler\n", func() {
			var n int64
			stream2, _ := os.Create("tmp/stream2.log")
			q, _ := Download(GET, "http://httpbin.org/bytes/10240", stream2)
			q.
				Progress(func(c, t, e int64) {
				n = c
				log.Println(c, t, e)
				So(c, ShouldBeLessThanOrEqualTo, t)
			}).Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
				l := int64(data.Len())
				So(n, ShouldEqual, l)
				So(res.ContentLength, ShouldEqual, l)
			})
		})
	})
}

func TestUpload(t *testing.T) {
	Convey("Uploading file", t, func(m C) {
		m.Convey("Uploading one file\n", func() {
			data := map[string]interface{}{
				"stream": "tmp/stream.log",
			}
			q, _ := Upload(POST, "http://httpbin.org/post", data)
			q.
				Progress(func(c, t, e int64) {
				log.Println(c, t, e)
				m.So(c, ShouldBeLessThanOrEqualTo, t)
			}).Do()
		})
		Convey("Uploading multi files\n", func() {
			stream2, _ := os.Open("quest_test.go")
			stream3 := bytes.NewBufferString(`Hello Quest!`)
			data := map[string]interface{}{
				"stream1": "quest.go", // filepath or filename
				"stream2": stream2,    // *os.File
				"stream3": stream3,    // io.Reader, filename is fieldname `stream3`
			}

			q, _ := Upload(POST, "http://httpbin.org/post", data)
			q.
				Parameters(map[string]string{"foo": "bar", "bar": "foo"}).
				Progress(func(c, t, e int64) {
				log.Println(c, t, e)
				m.So(c, ShouldBeLessThanOrEqualTo, t)
			}).Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
				l := int64(data.Len())
				m.So(res.ContentLength, ShouldEqual, l)
			})
		})
	})
}

func TestAuthenticate(t *testing.T) {
	type Auth struct {
		User          string
		Passwd        string
		Authenticated bool
	}
	user := "user"
	passwd := "password"

	Convey("Authenticate", t, func() {
		Convey("Basic Auth", func() {
			q, _ := Request(GET, "https://httpbin.org/basic-auth/"+user+"/"+passwd)
			q.Authenticate(user, passwd).
				ResponseJSON(func(_ *http.Request, _ *http.Response, data Auth, _ error) {
				So(data.User, ShouldEqual, user)
				So(data.Authenticated, ShouldEqual, true)
			}).Do()
		})
	})
}

func TestTimeout(t *testing.T) {
	Convey("Timeout", t, func() {
		Convey("It's timeout.", func() {
			s := time.Duration(3 * time.Second)
			q, _ := Request(GET, "https://httpbin.org/delay/5")
			q.Timeout(s).Do()
		})
		Convey("It's not timeout.", func() {
			s := time.Duration(30 * time.Second)
			q, _ := Request(GET, "https://httpbin.org/delay/5")
			q.Timeout(s).Do()
		})
	})
}

func TestSetHeader(t *testing.T) {
	type DataStruct struct {
		Headers map[string]string
	}
	Convey("set header", t, func() {
		q, _ := Request(GET, "http://httpbin.org/headers")
		q.Set("Quest", "Test").
			ResponseJSON(func(_ *http.Request, _ *http.Response, data DataStruct, _ error) {
			So(data.Headers["Quest"], ShouldEqual, "Test")
		}).Do()
	})
}

func TestBasicAuth(t *testing.T) {
	type DataStruct struct {
		Headers map[string]string
	}
	Convey("set header", t, func() {
		q, _ := Request(GET, "http://httpbin.org/headers")
		q.BasicAuth("test", "1234").
			ResponseJSON(func(_ *http.Request, _ *http.Response, data DataStruct, _ error) {
			username, password, ok := parseBasicAuth(data.Headers["Authorization"])
			So(username, ShouldEqual, "test")
			So(password, ShouldEqual, "1234")
			So(ok, ShouldEqual, true)
		}).Do()
	})
}

func TestBytesNotHandler(t *testing.T) {
	queryParams := url.Values{}
	queryParams.Set("foo", "bar")
	queryParams.Set("name", "活力")

	Convey("Response Bytes not handler", t, func() {
		q, _ := Request(GET, "http://httpbin.org/get")
		_, err := q.Query(&queryParams).Bytes()
		So(err, ShouldBeNil)
	})
}

func TestStringNotHandler(t *testing.T) {
	Convey("Response String not handler", t, func() {
		q, _ := Request(GET, "http://httpbin.org/get")
		_, err := q.String()
		So(err, ShouldBeNil)
	})
}

func TestJSONNotHandler(t *testing.T) {
	type DataStruct struct {
		Headers map[string]string
		Origin  string
	}
	Convey("Response JSON not handler", t, func() {
		q, _ := Request(GET, "http://httpbin.org/get")
		b := DataStruct{}
		err := q.JSON(&b)
		So(err, ShouldBeNil)
		So(b.Headers["Host"], ShouldEqual, "httpbin.org")
	})
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	if !strings.HasPrefix(auth, "Basic ") {
		return
	}
	c, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
