package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	. "github.com/go-libs/methods"
	"github.com/go-libs/progress"
	"github.com/go-libs/quest/utils"
	"github.com/go-libs/syncreader"
)

type HandlerFunc func(*http.Request, *http.Response, *bytes.Buffer, error)
type BytesHandlerFunc func(*http.Request, *http.Response, []byte, error)
type StringHandlerFunc func(*http.Request, *http.Response, string, error)

// A Request manages communication with http service.
type Request struct {
	// HTTP method
	Method Method

	// Base URL for Requests.
	Endpoint string
	Url      *url.URL

	// HTTP request
	req *http.Request

	// HTTP response
	res *http.Response

	// HTTP client
	client *http.Client

	// request header & body
	Header  http.Header
	Body    io.ReadCloser
	Length  int64
	rawBody interface{}

	isBodyClosed bool
	Buffer       *bytes.Buffer

	err error

	// Upload
	IsUpload bool
	files    map[string]interface{}
	fields   map[string]string

	// Download
	IsDownload  bool
	destination string

	// Progress
	pg *progress.Progress
}

func (r *Request) Files(files map[string]interface{}) *Request {
	r.files = files
	return r
}

func (r *Request) Fields(fields map[string]string) *Request {
	r.fields = fields
	return r
}

func (r *Request) Destionation(destination string) *Request {
	r.destination = destination
	return r
}

func (r *Request) Query(data interface{}) *Request {
	qs, err := utils.QueryString(data)
	r.err = err
	r.Url.RawQuery = qs
	return r
}

func (r *Request) Parameters(data interface{}) *Request {
	if encodesParametersInURL(r.Method) {
		return r
	}
	r.rawBody = data
	return r
}

func (r *Request) packBody() {
	if r.rawBody == nil {
		return
	}
	body, length, err := utils.PackBody(r.rawBody)
	r.err = err
	if length > 0 && body != nil {
		r.Body = body
		r.Length = length
	}
}

func (r *Request) Form(files map[string]interface{}, fields map[string]string) *Request {
	var data interface{}
	if len(files) > 0 {
		var (
			c  = make(chan bool, 1)
			b  = new(bytes.Buffer)
			mw = multipart.NewWriter(b)
		)
		//pr, pw := io.Pipe()
		//pb := &ProgressBar{Progress: func(c, t, e int64) {}}
		//w := io.MultiWriter(pb, pw)
		//writer := multipart.NewWriter(w)
		go func() {
			var fp io.Writer
			for fieldname, file := range files {
				fh, name, err := utils.GetFile(file)
				if fh != nil && err == nil {
					file = fh
					fp, err = mw.CreateFormFile(fieldname, filepath.Base(name))
				} else {
					fp, err = mw.CreateFormField(fieldname)
				}
				if err != nil {
					log.Fatal(err)
				}
				_, err = io.Copy(fp, ioutil.NopCloser(file.(io.Reader)))
				if err != nil {
					log.Fatal(err)
				}
			}

			for k, v := range fields {
				mw.WriteField(k, v)
			}
			mw.Close()
			//pw.Close()
			c <- true
		}()
		<-c
		//ppr := progress.NewReader()
		////ppr.Total = <-cc
		//ppr.Progress = func(c, t, e int64) {
		//	log.Println("Uploading stream", c, t, e)
		//}
		//rr := syncreader.New(pr, ppr)
		//log.Println(buff.Len())
		//body, length := PackBodyByReader(pr)
		r.Header.Set("Content-Type", mw.FormDataContentType())
		data = b
	} else {
		data = fields
	}
	r.rawBody = data
	return r
}

func (r *Request) Encoding(t string) *Request {
	t = strings.ToUpper(t)
	if t == "JSON" {
		t = "application/json"
	}
	if t != "" {
		r.Header.Set("Content-Type", t)
	}
	return r
}

func (r *Request) Authenticate(username, password string) *Request {
	return r
}

func (r *Request) Progress(f progress.HandlerFunc) *Request {
	r.pg = progress.New()
	r.pg.Progress = f
	return r
}

func (r *Request) response() (*bytes.Buffer, error) {
	if r.err != nil {
		return r.Buffer, r.err
	}
	if r.isBodyClosed {
		return r.Buffer, nil
	}
	r.isBodyClosed = true
	return r.Do()
}

func (r *Request) Response(handler HandlerFunc) *Request {
	_, err := r.response()
	handler(r.req, r.res, r.Buffer, err)
	return r
}

func (r *Request) ResponseBytes(handler BytesHandlerFunc) *Request {
	_, err := r.response()
	handler(r.req, r.res, r.Buffer.Bytes(), err)
	return r
}

func (r *Request) ResponseString(handler StringHandlerFunc) *Request {
	_, err := r.response()
	handler(r.req, r.res, r.Buffer.String(), err)
	return r
}

func (r *Request) ResponseJSON(f interface{}) *Request {
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

func (r *Request) Validate() *Request {
	return r
}

func (r *Request) validateAcceptContentType(map[string]string) bool {
	return true
}

// Acceptable Content Type
func (r *Request) ValidateAcceptContentType(map[string]string) bool {
	return true
}

func (r *Request) validateStatusCode(statusCodes ...int) bool {
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
func (r *Request) ValidateStatusCode(statusCodes ...int) *Request {
	r.response()
	if !r.validateStatusCode(statusCodes...) {
		r.err = errors.New("http: invalid status code " + strconv.Itoa(r.res.StatusCode))
	}
	return r
}

func (r *Request) Cancel() {}

func (r *Request) Do() (*bytes.Buffer, error) {
	r.req = &http.Request{
		Method: r.Method.String(),
		URL:    r.Url,
		Header: r.Header,
	}

	// uploading
	var hasPackBody bool
	if r.IsUpload {
		var fields map[string]string
		if r.rawBody != nil {
			fields, _ = r.rawBody.(map[string]string)
		}
		r.Form(r.files, fields)
		hasPackBody = true
		r.packBody()
		if r.pg != nil {
			r.pg.Total = r.Length
			r.Body = ioutil.NopCloser(syncreader.New(r.Body, r.pg))
		}
	}

	// pack body
	if !hasPackBody {
		r.packBody()
	}

	if r.Body != nil {
		if r.req.Header.Get("Content-Type") == "" {
			r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		r.req.Body = r.Body
		if r.Length > 0 {
			r.req.ContentLength = r.Length
		}
	}

	r.client = &http.Client{}
	res, err := r.client.Do(r.req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	r.res = res
	r.Buffer = new(bytes.Buffer)
	dw := io.MultiWriter(r.Buffer)

	// downloading
	if r.IsDownload {
		p, err := filepath.Abs(r.destination)
		if err != nil {
			return nil, err
		}
		f, err := os.Create(p)
		defer f.Close()
		if err != nil {
			return nil, err
		}
		if r.pg != nil {
			r.pg.Total = res.ContentLength
		}
		dw = io.MultiWriter(dw, r.pg, f)
		if err != nil {
			return nil, err
		}
	}

	_, err = io.Copy(dw, res.Body)
	if err != nil {
		return nil, err
	}
	return r.Buffer, nil
}

func (r *Request) Pipe() {
}

func encodesParametersInURL(method Method) bool {
	switch method {
	case GET, HEAD, DELETE:
		return true
	default:
		return false
	}
}
