# quest

Elegant HTTP Networking in Go.


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
  Response(func(req *http.Request, res *http.Response, data interface{}, e error) {
  fmt.Println(req)
  fmt.Println(res)
  fmt.Println(data)
  fmt.Println(err)
})
```


### Response Serialization

Built-in Response Methods

* `Response(*http.Request, *http.Response, interface{}, error)`
* `ResponseBytes(*http.Request, *http.Response, []byte, error)`
* `ResponseString(*http.Request, *http.Response, string, error)`
* `ResponseJSON(f interface{})`


#### Response String Handler

```go
quest.Request(GET, "http://httpbin.org/get").
  ResponseString(func(req *http.Request, res *http.Response, data string, e error) {
  fmt.Println(data)
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
  fmt.Println(data)
})

quest.Request(GET, "http://httpbin.org/get").
  ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
  fmt.Println(data)
})
```


#### Chained Response Handlers

Response handlers can even be chained:
```go
quest.Request(GET, "http://httpbin.org/get").
  ResponseString(func(req *http.Request, res *http.Response, data string, e error) {
  fmt.Println(data)
}).
  ResponseJSON(func(req *http.Request, res *http.Response, data *DataStruct, e error) {
  fmt.Println(data)
})
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
  fmt.Println(data)
})
  ResponseJSON(func(req *http.Request, res *http.Response, data OtherDataStruct, e error) {
  fmt.Println(data)
})
```
