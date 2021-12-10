package main

import (
	"context"
	"encoding/binary"
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

func receive(udpConn net.UDPConn) {
	log.Printf("Starting UDP ping receive loop")

	var nextSeq uint8
	var lost int
	for {

		p := &probe{}

		if err := binary.Read(&udpConn, binary.BigEndian, p); err != nil {
			return
		}

		log.Printf("Received probe %d", p.SeqNum)
		if p.SeqNum < nextSeq {
			log.Printf("Out of order packet seq/expected: %d/%d", p.SeqNum, nextSeq)
			lost -= 1
		} else if p.SeqNum > nextSeq {
			log.Printf("Out of order packet seq/expected: %d/%d", p.SeqNum, nextSeq)
			lost += int(p.SeqNum - nextSeq)
			nextSeq = p.SeqNum
		}

		latency := time.Now().UnixMilli() - p.SendTS
		log.Printf("E2E latency: %d ms", latency)
		log.Printf("Lost packets: %d", lost)
		nextSeq++
	}
}

func main() {
	server := flag.String("server", "127.0.0.1", "UDP server IP")
	port := flag.Int("port", 32767, "UDP server port")
	flag.Parse()

	rAddr := &net.UDPAddr{
		IP:   net.ParseIP(*server),
		Port: *port,
	}

	udpConn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		log.Fatalf("failed to DialUDP: %s", err)
	}
	defer udpConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	setupSigHandlers(cancel)

	ticker := time.NewTicker(probeInterval)

	log.Printf("Starting UDP ping client")
	go receive(*udpConn)

	var seq uint8
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down UDP client")
			return
		case <-ticker.C:
			log.Printf("Sending probe %d", seq)
			p := &probe{
				SeqNum: seq,
				SendTS: time.Now().UnixMilli(),
			}

			if err := udpConn.SetWriteDeadline(time.Now().Add(retryTimeout)); err != nil {
				log.Fatalf("failed to SetWriteDeadline: %s", err)
			}

			if err := binary.Write(udpConn, binary.BigEndian, p); err != nil {
				log.Fatalf("failed to binary.Write: %v", err)
			}

			seq++
		}
	}

}
