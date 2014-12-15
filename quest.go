package quest

import (
	"net/http"
	stdurl "net/url"
	"strings"
)

// base request client
func Request(method, endpoint string) (r *Requester, err error) {
	url, err := stdurl.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	r = &Requester{
		Method:   strings.ToUpper(method),
		Endpoint: endpoint,
		Url:      url,
		Header:   make(http.Header),
		timeout:  defaultTimeout,
	}
	return
}

// upload file / data / stream
func Upload(method, endpoint string, files map[string]interface{}) (r *Requester, err error) {
	r, err = Request(method, endpoint)
	if err != nil {
		return
	}
	r.IsUpload = true
	r.Files(files)
	return
}

// download file / data / stream to file
func Download(method, endpoint string, destination interface{}) (r *Requester, err error) {
	r, err = Request(method, endpoint)
	if err != nil {
		return
	}
	r.IsDownload = true
	r.Destination(destination)
	return
}
