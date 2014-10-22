package quest

import (
	"bytes"
	"encoding/json"
	"strconv"

	//"io"
	//"io/ioutil"
	"errors"
	"net/http"
	"net/url"

	. "github.com/go-libs/methods"
)

type JSONMaps map[string]interface{}

type HandlerFunc func(*http.Request, *http.Response, interface{}, error)
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

	err error
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
	if r.err != nil {
		return r.Buffer, r.err
	}
	if r.isBodyClosed {
		return r.Buffer, nil
	}
	r.isBodyClosed = true
	return r.Do()
}

func (r *Qrequest) Response(handler HandlerFunc) *Qrequest {
	body, err := r.response()
	handler(r.req, r.res, body, err)
	return r
}

func (r *Qrequest) ResponseBytes(handler BytesHandlerFunc) *Qrequest {
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
	if err != nil {
		handler(r.req, r.res, nil, err)
	} else {
		data := JSONMaps{}
		err = json.Unmarshal(body.Bytes(), &data)
		handler(r.req, r.res, data, err)
	}
	return r
}

func (r *Qrequest) Validate() *Qrequest {
	return r
}

func (r *Qrequest) validateStatusCode(statusCodes ...int) bool {
	statusCode := r.res.StatusCode
	if len(statusCodes) > 0 {
		for _, c := range statusCodes {
			if statusCode == c {
				return true
			}
		}
		// 200 <= x < 300
	} else if statusCode >= 200 && statusCode < 300 {
		return true
	}
	return false
}

func (r *Qrequest) ValidateStatusCode(statusCodes ...int) *Qrequest {
	r.response()
	if !r.validateStatusCode(statusCodes...) {
		r.err = errors.New("http: invalid status code " + strconv.Itoa(r.res.StatusCode))
	}
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
