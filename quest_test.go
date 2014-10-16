package quest

import (
	"net/http"
	"testing"

	. "github.com/go-libs/methods"
)

func TestString(t *testing.T) {
	Request(GET, "http://httpbin.org/get").
		Response(func(request *http.Request, response *http.Response, data interface{}, err error) {
	})
}
