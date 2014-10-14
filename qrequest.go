package quest

import (
	"net/http"
	"net/url"

	. "github.com/go-libs/methods"
)

type Qrequest struct {
	req *http.Request
	res http.Response
}

func (r *Qrequest) init() *Qrequest {
	return r
}

func (r *Qrequest) Query() *Qrequest {
	return r
}

func (r *Qrequest) Parameters() *Qrequest {
	return r
}

func (r *Qrequest) Authenticate(username, password string) *Qrequest {
	return r
}

func (r *Qrequest) Progress() *Qrequest {
	return r
}

func (r *Qrequest) Response() *Qrequest {
	return r
}

func (r *Qrequest) ResponseString() *Qrequest {
	return r
}

func (r *Qrequest) ResponseJSON() *Qrequest {
	return r
}

func (r *Qrequest) Validate() *Qrequest {
	return r
}

func (r *Qrequest) Cancel() {}

// Helpers:
func encodesParametersInURL(method Method) bool {
	switch method {
	case GET, HEAD, DELETE:
		return true
	default:
		return false
	}
}

func escape(s string) string {
	return url.QueryEscape(s)
}
