package quest

import (
	"net/http"
	stdurl "net/url"

	. "github.com/go-libs/methods"
)

func Request(method Method, endpoint string) (q *Qrequest) {
	url, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		panic(err)
	}
	q = &Qrequest{
		Method:   method,
		Endpoint: endpoint,
		Url:      url,
		Header:   make(http.Header),
	}
	return
}

// Upload File / Data / Stream
func Upload(method Method, endpoint string, files, fields map[string]string) (q *Qrequest) {
	q = Request(method, endpoint)
	q.isUpload = true
	q.files = files
	q.fields = fields
	return
}

// Download File / Data / Stream to file
func Download(method Method, endpoint, destination string) (q *Qrequest) {
	q = Request(method, endpoint)
	q.isDownload = true
	q.destination = destination
	return
}

func Println() {}

func DebugPrintln() {}
