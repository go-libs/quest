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
parameters := map[string]interface{}{
  "foo": []int{1, 2, 3},
  "bar": map[string]string{"baz": "qux"},
}

quest.Request(POST, "http://httpbin.org/post").
  Encoding("JSON").
  Parameters(parameters).
  ResponseJSON(func(request *http.Request, response *http.Response, data quest.JSONMaps, err error) {
  fmt.Printf("%+v\n", data["data"])
})
```
