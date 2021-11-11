package main

import (
	"fmt"
	"net/http"
)

func run() error {
	resp, err := http.Get("http://httpstat.us/304")
	if err != nil {
		return fmt.Errorf("Could not connect: %v", err)
	}

	switch {
	case resp.StatusCode >= 600:
		fmt.Println("Unknown")
	case resp.StatusCode >= 500:
		fmt.Println("Server Error")
	case resp.StatusCode >= 400:
		fmt.Println("Client Error")
	case resp.StatusCode >= 300:
		fmt.Println("Redirect")
	case resp.StatusCode >= 200:
		fmt.Println("Success")
	case resp.StatusCode >= 100:
		fmt.Println("Informational")
	default:
		fmt.Println("Incorrect")
	}
	return nil
}

func main() {

	fmt.Println(run())

}
