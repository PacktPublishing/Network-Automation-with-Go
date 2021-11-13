package main

import (
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	defer fmt.Println("Deferred cleanup")

	fmt.Println("Response status:", resp.Status)
}
