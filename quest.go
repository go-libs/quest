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

func Upload(method Method, endpoint string, files, fields map[string]string) *Qrequest {
	r := Request(method, endpoint)
	r.isUpload = true
	r.files = files
	r.fields = fields
	return r
}

func Download(method Method, endpoint, destination string) *Qrequest {
	r := Request(method, endpoint)
	r.isDownload = true
	r.destination = destination
	return r
}

func Println() {}

func DebugPrintln() {}
