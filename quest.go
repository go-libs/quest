package quest

import (
	"net/http"
	stdurl "net/url"

	. "github.com/go-libs/methods"
	req "github.com/go-libs/quest/request"
)

func Request(method Method, endpoint string) (r *req.Request) {
	url, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		panic(err)
	}
	r = &req.Request{
		Method:   method,
		Endpoint: endpoint,
		Url:      url,
		Header:   make(http.Header),
	}
	return
}

// upload file / data / stream
func Upload(method Method, endpoint string, files, fields map[string]string) (r *req.Request) {
	r = Request(method, endpoint)
	r.IsUpload = true
	r.Files(files)
	r.Fields(fields)
	return
}

// download file / data / stream to file
func Download(method Method, endpoint, destination string) (r *req.Request) {
	r = Request(method, endpoint)
	r.IsDownload = true
	r.Destionation(destination)
	return
}

func Println() {}

func DebugPrintln() {}
