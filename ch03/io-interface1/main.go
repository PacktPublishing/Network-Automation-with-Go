package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	src := strings.NewReader("The text")
	dst, err := os.Create("./file.txt")
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	io.Copy(dst, src)
}
