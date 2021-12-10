package main

import (
	"context"
	"flag"
	"log"
	"net"
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

func main() {
	port := flag.Int("port", listenPort, "UDP listen port")
	flag.Parse()

	listenSoc := &net.UDPAddr{
		IP:   net.ParseIP(listenAddr),
		Port: *port,
	}

	udpConn, err := net.ListenUDP("udp", listenSoc)
	if err != nil {
		log.Fatalf("failed to listen on %s:%d: %s", listenAddr, *port, err)
	}
	defer udpConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	setupSigHandlers(cancel)

	if err = udpConn.SetReadBuffer(maxReadBuffer); err != nil {
		log.Fatalf("failed to SetReadBuffer: %s", err)
	}

	log.Printf("Starting the UDP ping server")
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down UDP server")
			return
		default:
			bytes := make([]byte, maxReadBuffer)

			if err := udpConn.SetReadDeadline(time.Now().Add(retryTimeout)); err != nil {
				log.Fatalf("failed to SetReadDeadline: %s", err)
			}

			len, raddr, err := udpConn.ReadFromUDP(bytes)
			if err != nil {
				log.Printf("failed to ReadFromUDP: %s", err)
				continue
			}
			log.Printf("Received a probe from %s:%d", raddr.IP.String(), raddr.Port)

			if len == 0 {
				log.Printf("Received packet with 0 length")
				continue
			}

			n, err := udpConn.WriteToUDP(bytes[:len], raddr)
			if err != nil {
				log.Fatalf("Failed to WriteToUDP: %s", err)
			}

			if n != len {
				log.Printf("could not send the full packet")
			}
		}
	}

}
