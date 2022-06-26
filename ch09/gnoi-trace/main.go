package main

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/openconfig/gnoi/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	target       = "clab-netgo-ceos:6030"
	username     = "admin"
	password     = "admin"
	source       = "203.0.113.3"
	destinations = []string{
		"203.0.113.251",
		"203.0.113.252",
		"203.0.113.253",
	}
)

func main() {
	log.Println("Checking if routes have different paths")
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	sysSvc := system.NewSystemClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"username",
		username,
		"password",
		password,
	)

	var wg sync.WaitGroup
	wg.Add(len(destinations))
	traceCh := make(chan map[string][]mapset.Set, len(destinations))

	for _, dest := range destinations {
		go func(d string) {
			defer wg.Done()

			retryMax := 3
			retryCount := 0

		START:
			response, err := sysSvc.Traceroute(ctx, &system.TracerouteRequest{
				Destination: d,
				Source:      source,
			})
			if err != nil {
				log.Fatalf("Cannot trace path: %v", err)
			}
			var route []mapset.Set
			for {
				resp, err := response.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					log.Fatalf("Cannot receive response: %v", err)
				}
				if resp.Address == "" {
					continue
				}

				if int(resp.Hop) > len(route)+1 {
					log.Printf("Missed at least one hop in %s", d)
					if retryCount > retryMax-1 {
						log.Printf("reached retryMax, aborting %s", d)
						goto FINISH
					}
					log.Printf("retrying %s", d)
					retryCount += 1
					goto START
				}

				if len(route) < int(resp.Hop) {
					route = append(route, mapset.NewSet())
				}
				route[resp.Hop-1].Add(resp.Address)
			}

		FINISH:
			//log.Printf("Collected responses: %s, %+v", d, route)
			traceCh <- map[string][]mapset.Set{
				d: route,
			}

		}(dest)
	}
	wg.Wait()
	close(traceCh)

	routes := make(map[int]map[string]mapset.Set)

	for trace := range traceCh {
		for dest, paths := range trace {
			for hop, path := range paths {
				if _, ok := routes[hop]; !ok {
					routes[hop] = make(map[string]mapset.Set)
				}
				routes[hop][dest] = path
			}
		}
	}

	for hop, route := range routes {
		if hop == len(routes)-1 {
			continue
		}
		found := make(map[string]string)
		for myDest, myPaths := range route {
			for otherDest, otherPaths := range route {
				if myDest == otherDest {
					continue
				}
				diff := myPaths.Difference(otherPaths)
				if diff.Cardinality() == 0 {
					continue
				}

				v, ok := found[myDest]
				if ok && v == otherDest {
					continue
				}

				log.Printf("Found different paths at hop %d", hop)
				log.Printf("Destination %s: %+v", myDest, myPaths)
				log.Printf(
					"Destination %s: %+v",
					otherDest,
					otherPaths,
				)
				found[otherDest] = myDest
			}
		}
	}
	log.Println("Check complete")
}
