package quest

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/go-libs/methods"
	goquery "github.com/google/go-querystring/query"
)

type HandlerFunc func(*http.Request, *http.Response, *bytes.Buffer, error)
type BytesHandlerFunc func(*http.Request, *http.Response, []byte, error)
type StringHandlerFunc func(*http.Request, *http.Response, string, error)

type Qrequest struct {
	Method   Method
	Endpoint string
	Url      *url.URL
	req      *http.Request
	res      *http.Response
	client   *http.Client

	// request header & body
	Header http.Header
	Body   io.ReadCloser
	Length int64

	isBodyClosed bool
	Buffer       *bytes.Buffer

	err error

	Destination string
	isDownload  bool

	DataProgress func(current, total, expected int64)
}

func (r *Qrequest) QueryParameters(data interface{}) *Qrequest {
	var queryString string
	switch t := data.(type) {
	case string:
		queryString = t
		break
	case []byte:
		queryString = string(t)
		break
	case *url.Values:
		queryString = t.Encode()
		break
	default:
		qs, _ := goquery.Values(t)
		queryString = qs.Encode()
	}
	r.Url.RawQuery = queryString
	return r
}

func (r *Qrequest) Parameters(data interface{}) *Qrequest {
	if encodesParametersInURL(r.Method) {
		return r
	}
	var (
		body   io.ReadCloser
		length int64
	)
	switch t := data.(type) {
	case string:
		body, length = packBodyByString(t)
		break
	case []byte:
		body, length = packBodyByBytes(t)
		break
	case *url.Values:
		body, length = packBodyByString(t.Encode())
		break
	case *bytes.Buffer:
		body, length = packBodyByBytesBuffer(t)
		break
	case *bytes.Reader:
		body, length = packBodyByBytesReader(t)
		break
	case *strings.Reader:
		body, length = packBodyByStringsReader(t)
		break
	// JSON Object
	default:
		b, err := json.Marshal(data)
		if err != nil {
			r.err = err
			return r
		}
		body, length = packBodyByBytes(b)
	}
	if length > 0 && body != nil {
		r.Body = body
		r.Length = length
	}
	return r
}

func (r *Qrequest) Encoding(t string) *Qrequest {
	t = strings.ToUpper(t)
	if t == "JSON" {
		t = "application/json"
	}
	if t != "" {
		r.Header.Set("Content-Type", t)
	}
	return r
}

func (r *Qrequest) Authenticate(username, password string) *Qrequest {
	return r
}

func (r *Qrequest) Progress(f func(current, total, expected int64)) *Qrequest {
	r.DataProgress = f
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
	_, err := r.response()
	handler(r.req, r.res, r.Buffer, err)
	return r
}

func (r *Qrequest) ResponseBytes(handler BytesHandlerFunc) *Qrequest {
	_, err := r.response()
	handler(r.req, r.res, r.Buffer.Bytes(), err)
	return r
}

func (r *Qrequest) ResponseString(handler StringHandlerFunc) *Qrequest {
	_, err := r.response()
	handler(r.req, r.res, r.Buffer.String(), err)
	return r
}

func (r *Qrequest) ResponseJSON(f interface{}) *Qrequest {
	var (
		fn                    = reflect.ValueOf(f)
		t                     = reflect.TypeOf(f)
		argsNum               = t.NumIn()
		in                    = make([]reflect.Value, argsNum) //Panic if t is not kind of Func
		reqV, resV, dataV, eV reflect.Value
		err                   error
	)
	if argsNum != 4 {
		err = errors.New("ResponseJSON: invalid arguments.")
		return r
	} else {
		_, err = r.response()
		if err != nil {
			dataV = reflect.New(t.In(2)).Elem()
		} else {
			dataT := t.In(2)
			dataK := dataT.Kind()
			if dataK == reflect.Ptr {
				dataT = dataT.Elem()
			}
			dataN := reflect.New(dataT)
			data := dataN.Interface()
			err = json.Unmarshal(r.Buffer.Bytes(), &data)
			dataV = reflect.ValueOf(data)
			if dataK != reflect.Ptr {
				dataV = reflect.Indirect(dataV)
			}
		}
	}
	if err == nil {
		eV = reflect.New(t.In(3)).Elem()
	} else {
		eV = reflect.ValueOf(err)
	}
	reqV = reflect.ValueOf(r.req)
	resV = reflect.ValueOf(r.res)
	in[0] = reqV
	in[1] = resV
	in[2] = dataV
	in[3] = eV
	fn.Call(in)
	return r
}

func (r *Qrequest) Validate() *Qrequest {
	return r
}

func (r *Qrequest) validateAcceptContentType(map[string]string) bool {
	return true
}

// Acceptable Content Type
func (r *Qrequest) ValidateAcceptContentType(map[string]string) bool {
	return true
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

// Status Code
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
		URL:    r.Url,
		Header: r.Header,
	}
	if r.Length > 0 && r.Body != nil {
		if r.req.Header.Get("Content-Type") == "" {
			r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		r.req.Body = r.Body
		r.req.ContentLength = r.Length
	}
	r.client = &http.Client{}
	res, err := r.client.Do(r.req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r.res = res
	r.Buffer = new(bytes.Buffer)
	pb := &ProgressBar{
		Qreq:     r,
		Total:    res.ContentLength,
		Progress: r.DataProgress,
	}
	w := io.MultiWriter(r.Buffer, pb)

	if r.isDownload {
		p, err := filepath.Abs(r.Destination)
		if err != nil {
			return nil, err
		}
		f, err := os.Create(p)
		if err != nil {
			return nil, err
		}
		w = io.MultiWriter(w, f)
		defer f.Close()
		if err != nil {
			return nil, err
		}
	}
	_, err = io.Copy(w, res.Body)
	if err != nil {
		return nil, err
	}
	return r.Buffer, nil
}

// Helpers:

func packBodyByString(s string) (io.ReadCloser, int64) {
	return ioutil.NopCloser(bytes.NewBufferString(s)), int64(len(s))
}

func packBodyByBytes(b []byte) (io.ReadCloser, int64) {
	return ioutil.NopCloser(bytes.NewBuffer(b)), int64(len(b))
}

func packBodyByBytesBuffer(b *bytes.Buffer) (io.ReadCloser, int64) {
	return ioutil.NopCloser(b), int64(b.Len())
}

func packBodyByBytesReader(b *bytes.Reader) (io.ReadCloser, int64) {
	return ioutil.NopCloser(b), int64(b.Len())
}

func packBodyByStringsReader(b *strings.Reader) (io.ReadCloser, int64) {
	return ioutil.NopCloser(b), int64(b.Len())
}

func encodesParametersInURL(method Method) bool {
	switch method {
	case GET, HEAD, DELETE:
		return true
	default:
		return false
	}
}
