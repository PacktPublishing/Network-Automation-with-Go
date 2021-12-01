package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	listenAddr     = "0.0.0.0"
	listenPort     = 32767
	probeSizeBytes = 9
	maxReadBuffer  = 425984
	retryTimeout   = time.Second * 5
	probeInterval  = time.Second
)

func setupSigHandlers(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-sigs
		log.Printf("Received syscall: %+v", sig)
		cancel()
	}()

}

type probe struct {
	SeqNum uint8
	SendTS int64
}
