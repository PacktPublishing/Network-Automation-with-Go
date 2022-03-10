package main

import (
	"flag"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"golang.org/x/net/bpf"
)

var (
	intf = flag.String("intf", "eth1", "interface")
)

func main() {
	handle, err := pcapgo.NewEthernetHandle(*intf)
	if err != nil {
		log.Panic(err)
	}

	rawInstructions, err := bpf.Assemble([]bpf.Instruction{
		// Load "EtherType" field from the ethernet header.
		bpf.LoadAbsolute{Off: 12, Size: 2},
		// Skip to the last instruction if EtherType is not IPv4.
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0x800, SkipTrue: 3},
		// Load "Protocol" field from the IPv4 header.
		bpf.LoadAbsolute{Off: 23, Size: 1},
		// Skip to the last instruction if Protocol is not UDP.
		bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 0x11, SkipTrue: 1},
		// Verdict is "send up to 4k of the packet to userspace."
		bpf.RetConstant{Val: 4096},
		// Verdict is "ignore packet."
		bpf.RetConstant{Val: 0},
	})
	if err != nil {
		log.Panic(err)
	}
	if err := handle.SetBPF(rawInstructions); err != nil {
		log.Panic(err)
	}

	packetSource := gopacket.NewPacketSource(
		handle,
		layers.LayerTypeEthernet,
	)
	for packet := range packetSource.Packets() {
		if l4 := packet.TransportLayer(); l4 == nil {
			continue
		}

		sflowLayer := packet.Layer(layers.LayerTypeSFlow)
		if sflowLayer != nil {
			sflow, ok := sflowLayer.(*layers.SFlowDatagram)
			if !ok {
				log.Println("failed decoding sflow")
				continue
			}

			for _, sample := range sflow.FlowSamples {
				for _, record := range sample.GetRecords() {
					p, ok := record.(layers.SFlowRawPacketFlowRecord)
					if !ok {
						log.Println("failed to decode sflow record")
						continue
					}

					srcIP, dstIP := p.Header.
						NetworkLayer().
						NetworkFlow().
						Endpoints()
					sPort, dPort := p.Header.
						TransportLayer().
						TransportFlow().
						Endpoints()
					log.Printf("flow record: %s:%s <-> %s:%s\n",
						srcIP,
						sPort,
						dstIP,
						dPort,
					)
				}

			}

		}

	}

}
