package main

import "fmt"

const F = 42

func main() {

	if Version != "" {
		fmt.Printf("Version: %q\n", Version)
	}

	if GitCommit != "" {
		fmt.Printf("Git Commit: %q\n", GitCommit)
	}

	fmt.Println("Hello World")

}
