package quest

// Method.
type Method int

// Methods.
//  HTTP method definitions.
//  See http://tools.ietf.org/html/rfc7231#section-4.3
const (
	OPTIONS Method = iota
	GET
	HEAD
	POST
	PUT
	PATCH
	DELETE
	TRACE
	CONNECT
)

// Method map.
var Methods = map[Method]string{
	OPTIONS: "OPTIONS",
	GET:     "GET",
	HEAD:    "HEAD",
	POST:    "POST",
	PUT:     "PUT",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	TRACE:   "TRACE",
	CONNECT: "CONNECT",
}

// String returns the method string.
func (m Method) String() string {
	return Methods[m]
}
