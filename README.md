# quest

Elegant HTTP Networking in Go.

__[Docs]()__


## Usage


### Making a Request

```go
import "github.com/go-libs/quest"
import . "github.com/go-libs/methods"

quest.Request(GET, "http://httpbin.org/get")
```


### Response Handling

```go
quest.Request(GET, "http://httpbin.org/get").
  Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, e error) {
  log.Println(req)
  log.Println(res)
  log.Println(data)
  log.Println(err)
})
```


### Response Serialization

Built-in Response Methods

* `Response(func(*http.Request, *http.Response, *bytes.Buffer, error))`
* `ResponseBytes(func(*http.Request, *http.Response, []byte, error))`
* `ResponseString(func(*http.Request, *http.Response, string, error))`
* `ResponseJSON(f interface{})`, `f` ___Must___ be `func`
    - `func(req *http.Request, res *http.Response, data *JSONStruct, e error)`
    - `func(req *http.Request, res *http.Response, data JSONStruct, e error)`


#### Response String Handler

```go
quest.Request(GET, "http://httpbin.org/get").
  ResponseString(func(req *http.Request, res *http.Response, data string, e error) {
  log.Println(data)
})
```


#### Response JSON Handler

```go
type DataStruct struct {
  Headers map[string]string
  Origin  string
}

quest.Request(GET, "http://httpbin.org/get").
  ResponseJSON(func(req *http.Request, res *http.Response, data DataStruct, e error) {
  log.Println(data)
})

quest.Request(GET, "http://httpbin.org/get").
  ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
  log.Println(data)
})
```


#### Chained Response Handlers

Response handlers can even be chained:
```go
quest.Request(GET, "http://httpbin.org/get").
  ResponseString(func(req *http.Request, res *http.Response, data string, e error) {
  log.Println(data)
}).
  ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
  log.Println(data)
})
```


### HTTP Methods

```go
import . "github.com/go-libs/methods"
```


### Query String

```go
type Options struct {
  Foo string `url:"foo"`
}

quest.Request(GET, "http://httpbin.org/get").
  QueryParameters(Options{"bar"})
// http://httpbin.org/get?foo=bar
```


### POST Request with JSON-encoded Parameters

```go
type PostParameters struct {
  Foo []int             `json:"foo,omitempty"`
  Bar map[string]string `json:"bar,omitempty"`
}

parameters := PostParameters{
  "foo": []int{1, 2, 3},
  "bar": map[string]string{"baz": "qux"},
}

type DataStruct struct {
  Headers map[string]string
  Origin  string
  Json    PostParameters `json:"json,omitempty"`
}

type OtherDataStruct struct {
  Headers map[string]string
  Origin  string
}

quest.Request(POST, "http://httpbin.org/post").
  Encoding("JSON").
  Parameters(&parameters).
  ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
  log.Println(data)
})
  ResponseJSON(func(req *http.Request, res *http.Response, data OtherDataStruct, e error) {
  log.Println(data)
})
```



### Downloading


#### Downloading a File

```go
quest.Download(GET, "http://httpbin.org/stream/100", "stream.log").Do()
```


#### Downloading a File w/Progress

```go
destination := "tmp/stream.log"
quest.Download(GET, "http://httpbin.org/bytes/1024", destination).
  Progress(func(bytesRead, totalBytesRead, totalBytesExpectedToRead int64) {
    log.Println(bytesRead, totalBytesRead, totalBytesExpectedToRead)
  }).Do()

destination := "tmp/stream2.log"
quest.Download(GET, "http://httpbin.org/bytes/10240", destination).
  Progress(func(current, total, expected int64) {
    log.Println(current, total, expected)
  }).Response(func(request *http.Request, response *http.Response, data *bytes.Buffer, err error) {
    log.Println(data.Len())
  })
```


### Uploading


#### Uploading a File

```go
quest.Upload(POST, "http://httpbin.org/post", map[string]string{"stream": "tmp/stream.log"}, nil).Do()
```


#### Uploading multi files and in progress

```go
quest.Upload(POST, "http://httpbin.org/post", map[string]string{"stream": "tmp/stream.log", "stream2": "tmp/stream2.log"}, nil).
  Progress(func(current, total, expected int64) {
    log.Println(current, total, expected)
  }).Response(func(req *http.Request, res *http.Response, data *bytes.Buffer, err error) {
    log.Println(data.Len())
  })
```


## License

MIT

[Docs]: http://godoc.org/github.com/go-libs/quest
