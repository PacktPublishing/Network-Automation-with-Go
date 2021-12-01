package main

import (
	"context"
	"encoding/binary"
	"flag"
	"log"
	"net"
	"time"
)

func receive(udpConn net.UDPConn) {
	log.Printf("Starting UDP ping receive loop")

	var nextSeq uint8
	for {

		p := &probe{}

		if err := binary.Read(&udpConn, binary.BigEndian, p); err != nil {
			return
		}

		if p.SeqNum < nextSeq {
			log.Printf("Out of order packet seq/expected: %d/%d", p.SeqNum, nextSeq)
		} else if p.SeqNum > nextSeq {
			log.Printf("Out of order packet seq/expected: %d/%d", p.SeqNum, nextSeq)
			nextSeq = p.SeqNum
		}
		log.Printf("Current TS: %d", time.Now().Unix())
		log.Printf("Received TS: %d", p.SendTS)
		latency := time.Unix(0, time.Now().Unix()-p.SendTS)
		log.Printf("E2E latency: %d ns", latency.Nanosecond())

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
				SendTS: time.Now().Unix(),
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
