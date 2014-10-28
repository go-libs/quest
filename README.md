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
