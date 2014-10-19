package quest

import (
	"bytes"
	"encoding/json"
	//"io"
	//"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/go-libs/methods"
)

type JSONMaps map[string]interface{}

type BytesHandlerFunc func(*http.Request, *http.Response, []byte, error)
type StringHandlerFunc func(*http.Request, *http.Response, string, error)
type JSONHandlerFunc func(*http.Request, *http.Response, JSONMaps, error)

type Qrequest struct {
	Method Method
	Url    string
	Uri    *url.URL
	req    *http.Request
	res    *http.Response
	client *http.Client

	isBodyClosed bool
	Buffer       *bytes.Buffer
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

func (r *Qrequest) response() (*bytes.Buffer, error) {
	if r.isBodyClosed {
		return r.Buffer, nil
	}
	r.isBodyClosed = true
	return r.Do()
}

func (r *Qrequest) Response(handler BytesHandlerFunc) *Qrequest {
	body, err := r.response()
	handler(r.req, r.res, body.Bytes(), err)
	return r
}

func (r *Qrequest) ResponseString(handler StringHandlerFunc) *Qrequest {
	body, err := r.response()
	handler(r.req, r.res, body.String(), err)
	return r
}

func (r *Qrequest) ResponseJSON(handler JSONHandlerFunc) *Qrequest {
	body, err := r.response()
	data := JSONMaps{}
	err = json.Unmarshal(body.Bytes(), &data)
	handler(r.req, r.res, data, err)
	return r
}

func (r *Qrequest) Validate() *Qrequest {
	return r
}

func (r *Qrequest) Cancel() {}

func (r *Qrequest) Do() (*bytes.Buffer, error) {
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
	defer res.Body.Close()
	r.Buffer = new(bytes.Buffer)
	r.Buffer.ReadFrom(res.Body)
	return r.Buffer, nil
}

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
