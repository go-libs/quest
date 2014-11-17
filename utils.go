package quest

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

func nopCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

func packBodyByString(s string) (io.ReadCloser, int64) {
	return nopCloser(bytes.NewBufferString(s)), int64(len(s))
}

func packBodyByBytes(b []byte) (io.ReadCloser, int64) {
	return nopCloser(bytes.NewBuffer(b)), int64(len(b))
}

func packBodyByBytesBuffer(b *bytes.Buffer) (io.ReadCloser, int64) {
	return nopCloser(b), int64(b.Len())
}

func packBodyByBytesReader(b *bytes.Reader) (io.ReadCloser, int64) {
	return nopCloser(b), int64(b.Len())
}

func packBodyByPipeReader(pr *io.PipeReader) (io.ReadCloser, int64) {
	b := new(bytes.Buffer)
	length, _ := b.ReadFrom(pr)
	return nopCloser(b), length
}

func packBodyByReader(pr io.Reader) (io.ReadCloser, int64) {
	b := new(bytes.Buffer)
	length, _ := b.ReadFrom(pr)
	return nopCloser(b), length
}

func packBodyByStringsReader(b *strings.Reader) (io.ReadCloser, int64) {
	return nopCloser(b), int64(b.Len())
}
