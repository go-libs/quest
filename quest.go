package quest

import (
	"net/http"
	stdurl "net/url"
)

// base request client
func Request(method Method, endpoint string) (r *Requester, err error) {
	url, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	r = &Requester{
		Method:   method,
		Endpoint: endpoint,
		Url:      url,
		Header:   make(http.Header),
		timeout:  defaultTimeout,
	}
	return
}

// upload file / data / stream
func Upload(method Method, endpoint string, files map[string]interface{}) (r *Requester, err error) {
	r, err = Request(method, endpoint)
	if err != nil {
		return
	}
	r.IsUpload = true
	r.Files(files)
	return
}

// download file / data / stream to file
func Download(method Method, endpoint string, destination interface{}) (r *Requester, err error) {
	r, err = Request(method, endpoint)
	if err != nil {
		return
	}
	r.IsDownload = true
	r.Destination(destination)
	return
}

func Println() {}

func DebugPrintln() {}
