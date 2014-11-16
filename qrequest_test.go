package quest

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	. "github.com/go-libs/methods"
)

type JSONMaps map[string]interface{}

func TestString(t *testing.T) {
	return
	Request(GET, "http://httpbin.org/get").
		ValidateStatusCode().
		Response(func(request *http.Request, response *http.Response, data *bytes.Buffer, err error) {
		fmt.Println("Response")
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	}).ResponseBytes(func(request *http.Request, response *http.Response, data []byte, err error) {
		fmt.Println("ResponseBytes")
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	}).ResponseString(func(request *http.Request, response *http.Response, data string, err error) {
		fmt.Println("ResponseString")
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Println(err)
	}).ResponseJSON(func(request *http.Request, response *http.Response, data JSONMaps, err error) {
		fmt.Println("ResponseJSON")
		fmt.Println(request)
		fmt.Println(response)
		fmt.Println(data)
		fmt.Printf("%+v\n", data)
		fmt.Printf("%+v\n", data["headers"])
		fmt.Println(err)
	})
}
