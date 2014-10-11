package quest

import (
	. "github.com/go-libs/methods"
	"net/http"
)

type Qrequest struct {
	req *http.Request
	res http.Response
}

func (r *Qrequest) init() *Qrequest {
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
