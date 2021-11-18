package main

import (
	"fmt"
	"time"
)

func repeat(d chan bool, c <-chan time.Time) {
	for {
		select {
		case <-d:
			return
		case t := <-c:
			fmt.Println("Run at", t.Local())
		}
	}
}

func main() {
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go repeat(done, ticker.C)

	time.Sleep(2100 * time.Millisecond)
	ticker.Stop()
	done <- true
}
