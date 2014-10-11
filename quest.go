package quest

import (
	. "github.com/go-libs/methods"
)

func Request(method Method, url string) *Qrequest {
	request := &Qrequest{}
	request.init()
	return request
}

func Upload() {}

func Download() {}

func Println() {}

func DebugPrintln() {}
