package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type myReader struct {
	src io.Reader
}

func (r *myReader) Read(buf []byte) (int, error) {
	tmp := make([]byte, len(buf))
	n, err := r.src.Read(tmp)
	copy(buf[:n], bytes.Title(tmp[:n]))
	return n, err
}

func NewMyReader(r io.Reader) io.Reader {
	return &myReader{src: r}
}

func main() {
	r1 := strings.NewReader("network automation with go")
	r2 := NewMyReader(r1)

	io.Copy(os.Stdout, r2)
}
