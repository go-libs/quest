package quest

import (
	"net/http"
	stdurl "net/url"
	"os"

	. "github.com/go-libs/methods"
)

func Request(method Method, endpoint string) *Qrequest {
	url, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	request := &Qrequest{
		Method:   method,
		Endpoint: endpoint,
		Url:      url,
		Header:   make(http.Header),
	}
	return request
}

func Upload(method Method, endpoint string) *Qrequest {
	return Request(method, endpoint)
}

func Download(method Method, endpoint string) *Qrequest {
	return Request(method, endpoint)
}

func Println() {}

func DebugPrintln() {}
