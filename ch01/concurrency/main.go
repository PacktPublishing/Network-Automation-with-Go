package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func get_config(d string) string {
	time.Sleep(time.Duration(rand.Intn(1)) * time.Second)
	return fmt.Sprintf("Connected to device %q", d)
}

func connect(devices []string, results chan string) {
	var wg sync.WaitGroup
	wg.Add(len(devices))

	for _, d := range devices {
		go func(d string) {
			defer wg.Done()
			results <- get_config(d)
		}(d)
	}

	wg.Wait()
	close(results)
}

func main() {

	devices := []string{"leaf01", "leaf02", "spine01"}
	resultCh := make(chan string, len(devices))

	go connect(devices, resultCh)

	fmt.Println("Continuing execution")

	for msg := range resultCh {
		fmt.Println(msg)
	}
}
