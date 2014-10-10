package quest

import (
  "net/http"
)

type Request struct {
  req *http.Request
  res http.Response
}

func (r *Request) Authenticate(username, password string) *Request {
  return r
}

func (r *Request) Progress() *Request {}

func (r *Request) Response() *Request {
  return r
}

func (r *Request) Validate() *Request {}

func (r *Request) Cancel() {}
