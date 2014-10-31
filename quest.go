package quest

import (
	"net/http"
	stdurl "net/url"
	"os"

	. "github.com/go-libs/methods"
)

func Request(method Method, url string) *Qrequest {
	uri, err := stdurl.ParseRequestURI(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	request := &Qrequest{
		Method: method,
		Url:    url,
		Uri:    uri,
		Header: http.Header{},
	}
	return request
}

func Upload(method Method, url string) *Qrequest {
	return Request(method, url)
}

func Download(method Method, url string) *Qrequest {
	return Request(method, url)
}

func Println() {}

func DebugPrintln() {}
