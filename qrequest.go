package quest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/go-libs/methods"
)

type HandlerFunc func(*http.Request, *http.Response, interface{}, error)

type Qrequest struct {
	Method Method
	Url    string
	Uri    *url.URL
	req    *http.Request
	res    *http.Response
	client *http.Client
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

func (r *Qrequest) response() (io.ReadCloser, error) {
	r.req = &http.Request{
		Method: r.Method.String(),
		URL:    r.Uri,
	}
	r.client = &http.Client{}
	res, err := r.client.Do(r.req)
	if err != nil {
		return nil, err
	}
	r.res = res
	r.Do()
	defer res.Body.Close()
	return res.Body, nil
}

func (r *Qrequest) Response(handler HandlerFunc) *Qrequest {
	body, err := r.response()
	handler(r.req, r.res, body.(io.ReadCloser), err)
	return r
}

func (r *Qrequest) ResponseString(handler HandlerFunc) *Qrequest {
	body, err := r.response()
	data, err := ioutil.ReadAll(body)
	handler(r.req, r.res, string(data), err)
	return r
}

func (r *Qrequest) ResponseJSON(handler HandlerFunc) *Qrequest {
	return r
}

func (r *Qrequest) Validate() *Qrequest {
	return r
}

func (r *Qrequest) Cancel() {}

func (r *Qrequest) Do() {}

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
