package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

type Router struct {
	Hostname  string `yaml:"hostname"`
	Platform  string `yaml:"platform"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	StrictKey bool   `yaml:"strictkey"`
}

type Inventory struct {
	Routers []Router `yaml:"router"`
}

var m sync.RWMutex = sync.RWMutex{}

func getVersion(r Router, out chan string, wg *sync.WaitGroup, isAlive map[string]bool) {
	defer wg.Done()
	
	//m.Lock()
	isAlive[r.Hostname] = true
	//m.Unlock()

	out <- "test"
}

func printer(in chan string) {
	for out := range in {
		fmt.Printf("MESSAGE: %s\n", out)
	}
}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	src, err := os.Open("input.yml")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	d := yaml.NewDecoder(src)

	var inv Inventory
	err = d.Decode(&inv)
	if err != nil {
		panic(err)
	}

	ch := make(chan string)

	isAlive := make(map[string]bool)

	go printer(ch)

	var wg sync.WaitGroup
	for _, v := range inv.Routers {
		wg.Add(1)
		go getVersion(v, ch, &wg, isAlive)
	}
	wg.Wait()

	m.RLock()
	for name, v := range isAlive {
		fmt.Printf("Router %s is alive: %t\n", name, v)
	}
	m.RUnlock()
	
	close(ch)
}