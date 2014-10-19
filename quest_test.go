package quest

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/go-libs/methods"
)

func TestString(t *testing.T) {
	Request(GET, "http://httpbin.org/get").
		Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	}).ResponseString(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	}).ResponseJSON(func(request *http.Request, response *http.Response, data interface{}, err error) {
		fmt.Println(request)
		fmt.Println(response)
		fmt.Printf("%+v\n", data)
		fmt.Println(err)
	})
}
