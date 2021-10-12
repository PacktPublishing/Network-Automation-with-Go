package main

import "fmt"

const F = 42

func main() {

	if Version != "" {
		fmt.Printf("Version: %s\n", Version)
	}

	if GitCommit != "" {
		fmt.Printf("Git Commit: %s\n", GitCommit)
	}

	fmt.Println("Hello World")

}
