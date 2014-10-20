package quest

import (
	. "github.com/go-libs/methods"
	stdurl "net/url"
	"os"
)

func Request(method Method, url string) *Qrequest {
	uri, err := stdurl.ParseRequestURI(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	request := &Qrequest{
		Method: method,
		Url:    url,
		Uri:    uri,
	}
	return request
}

func Upload() {}

func Download(method Method, url string) *Qrequest {
  return Request(method, url)
}

func Println() {}

func DebugPrintln() {}
