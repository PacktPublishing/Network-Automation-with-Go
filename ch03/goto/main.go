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
		goto exception
	case resp.StatusCode >= 500:
		fmt.Println("Server Error")
		goto failure
	case resp.StatusCode >= 400:
		fmt.Println("Client Error")
		goto failure
	case resp.StatusCode >= 300:
		fmt.Println("Redirect")
		goto exit
	case resp.StatusCode >= 200:
		fmt.Println("Success")
		goto exit
	case resp.StatusCode >= 100:
		fmt.Println("Informational")
		goto exit
	default:
		fmt.Println("Incorrect")
		goto exception
	}

	exception:
	panic("Unexpected response")

	failure:
	return fmt.Errorf("Failed to connect: %v", err)

	exit:
	fmt.Println("Connection successful")
	return nil
}

func main() {

	fmt.Println(run())

}
