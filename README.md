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
  Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
  fmt.Println(request)
  fmt.Println(response)
  fmt.Println(data)
  fmt.Println(err)
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
  ResponseJSON(func(request *http.Request, response *http.Response, data *DataStruct, err error) {
  fmt.Println(data)
})
  ResponseJSON(func(request *http.Request, response *http.Response, data OtherDataStruct, err error) {
  fmt.Println(data)
})
```
