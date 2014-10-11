package quest

import (
	. "github.com/go-libs/methods"
	"testing"
)

func TestString(t *testing.T) {
	Request(GET, "http://httpbin.org/get").Response()
}
