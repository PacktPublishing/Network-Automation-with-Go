package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	epb "github.com/cloudprober/cloudprober/probes/external/proto"
	"github.com/cloudprober/cloudprober/probes/external/serverutils"
	"github.com/jwhited/corebgp"
	bgp "github.com/osrg/gobgp/v3/pkg/packet/bgp"
	"google.golang.org/protobuf/proto"
)

const (
	bgpPingType = 42
	type42Len   = 46
	type42Flags = bgp.BGP_ATTR_FLAG_TRANSITIVE | bgp.BGP_ATTR_FLAG_OPTIONAL
)

var (
	id            = flag.String("id", "", "router ID")
	nlri          = flag.String("nlri", "", "local NLRI")
	localAS       = flag.Uint("las", 0, "local AS")
	remoteAS      = flag.Uint("ras", 0, "remote AS")
	localAddress  = flag.String("laddr", "", "local address")
	remoteAddress = flag.String("raddr", "", "remote address")
	passive       = flag.Bool("p", false, "passive mode")
	cloudprober   = flag.Bool("c", false, "cloudprober mode")
	interval      = flag.String("i", "10s", "probing interval")
)

type plugin struct {
	probe     net.IP
	host      []byte
	pingCh    chan ping
	probeCh   chan struct{}
	resultsCh chan string
	store     []string
	passive   bool
}

type ping struct {
	source []byte
	dest   []byte
	ts     int64
}

func (p *plugin) GetCapabilities(c corebgp.PeerConfig) []corebgp.Capability {
	caps := make([]corebgp.Capability, 0)
	caps = append(caps, newMPCap(1, 1))
	return caps
}

func (p *plugin) OnOpenMessage(peer corebgp.PeerConfig, routerID net.IP, capabilities []corebgp.Capability) *corebgp.Notification {
	return nil
}

func (p *plugin) OnEstablished(
	peer corebgp.PeerConfig,
	writer corebgp.UpdateMessageWriter,
) corebgp.UpdateMessageHandler {
	log.Println("peer established, starting main loop")
	go func() {
		// prepare withdraw timer
		withdrawAfter := time.Second * 7
		withdraw := time.NewTicker(withdrawAfter)
		withdraw.Stop()
		for {
			select {
			case pingReq := <-p.pingCh:
				src := string(bytes.Trim(pingReq.source, "\x00"))
				pingReq.dest = p.host
				type42PathAttr := bgp.NewPathAttributeUnknown(
					type42Flags,
					bgpPingType,
					buildPayload(pingReq),
				)
				bytes, err := p.buildUpdate(
					type42PathAttr,
					peer.LocalAddress,
					peer.LocalAS,
				)
				if err != nil {
					log.Fatal(err)
				}

				p.store = make([]string, 0)
				writer.WriteUpdate(bytes)
				withdraw.Reset(withdrawAfter)

				log.Printf("sent ping response to %s\n", src)

			case <-p.probeCh:
				log.Println("sending ping request")
				withdraw.Stop()

				pingReq := ping{
					source: p.host,
					ts:     time.Now().Unix(),
				}
				type42PathAttr := bgp.NewPathAttributeUnknown(
					type42Flags,
					bgpPingType,
					buildPayload(pingReq),
				)
				bytes, err := p.buildUpdate(
					type42PathAttr,
					peer.LocalAddress,
					peer.LocalAS,
				)
				if err != nil {
					log.Fatal(err)
				}

				// reset results
				p.store = make([]string, 0)
				writer.WriteUpdate(bytes)
				withdraw.Reset(withdrawAfter)

			case <-withdraw.C:
				log.Println("Sending ping withdraw")
				bytes, err := p.buildWithdraw()
				if err != nil {
					log.Fatal(err)
				}

				writer.WriteUpdate(bytes)
				withdraw.Stop()

				// optionally send results if receiver is ready
				select {
				case p.resultsCh <- strings.Join(p.store, "\n"):
				default:
				}

			}
		}
	}()

	return p.handleUpdate
}

func (p *plugin) OnClose(peer corebgp.PeerConfig) {
	log.Println("peer closed")
}

func (p *plugin) handleUpdate(
	peer corebgp.PeerConfig,
	update []byte,
) *corebgp.Notification {
	msg, err := bgp.ParseBGPBody(
		&bgp.BGPHeader{Type: bgp.BGP_MSG_UPDATE},
		update,
	)
	if err != nil {
		log.Fatal("failed to parse bgp message ", err)
	}

	if err := bgp.ValidateBGPMessage(msg); err != nil {
		log.Fatal("validate BGP message ", err)
	}

	for _, attr := range msg.Body.(*bgp.BGPUpdate).PathAttributes {
		// ignore all attributes except for 42
		if attr.GetType() != bgpPingType {
			continue
		}

		source, dest, ts, err := parseType42(attr)
		if err != nil {
			log.Fatal(err)
		}
		sourceHost := string(bytes.Trim(source, "\x00"))
		destHost := string(bytes.Trim(dest, "\x00"))

		// if source is us, it may be a response
		if sourceHost == *id {
			log.Println("Received a ping response")
			// some BGP stacks reflect all eBGP routes back to the sender
			if len(destHost) == 0 {
				log.Println("Received a looped response, ignoring")
				continue
			}

			// destHost is always set on a response
			rtt := time.Since(ts).Nanoseconds()
			metric := fmt.Sprintf(
				"bgp_ping_rtt_ms{device=%s} %f\n",
				destHost,
				float64(rtt)/1e6,
			)
			log.Println(metric)
			p.store = append(p.store, metric)
			return nil
		}

		// if source is not us and destHost is set, it must be an update from another responder
		if len(destHost) != 0 {
			continue
		}

		p.pingCh <- ping{source: source, ts: ts.Unix()}
		return nil
	}

	return nil
}

func (p *plugin) buildWithdraw() ([]byte, error) {
	myNLRI := bgp.NewIPAddrPrefix(32, p.probe.String())
	withdrawnRoutes := []*bgp.IPAddrPrefix{myNLRI}
	msg := bgp.NewBGPUpdateMessage(
		withdrawnRoutes,
		[]bgp.PathAttributeInterface{},
		nil,
	)
	return msg.Body.Serialize()
}

func (p *plugin) buildUpdate(
	type42 *bgp.PathAttributeUnknown,
	localAddr net.IP,
	localAS uint32,
) ([]byte, error) {
	withdrawnRoutes := []*bgp.IPAddrPrefix{}
	nexthop := bgp.NewPathAttributeNextHop(
		localAddr.String(),
	)
	asPath := bgp.NewPathAttributeAsPath(
		[]bgp.AsPathParamInterface{
			bgp.NewAs4PathParam(
				bgp.BGP_ASPATH_ATTR_TYPE_SEQ,
				[]uint32{localAS},
			),
		},
	)
	origin := bgp.NewPathAttributeOrigin(2) // origin incomplete
	myNLRI := bgp.NewIPAddrPrefix(32, p.probe.String())
	msg := bgp.NewBGPUpdateMessage(
		withdrawnRoutes,
		[]bgp.PathAttributeInterface{
			type42,
			nexthop,
			asPath,
			origin,
		},
		[]*bgp.IPAddrPrefix{myNLRI},
	)
	return msg.Body.Serialize()
}

func main() {
	flag.Parse()

	if *cloudprober && *passive {
		log.Fatal("can't use both 'cloudprober' and 'passive' mods")
	}

	timeInterval, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("starting corebgp server")
	srv, err := corebgp.NewServer(net.ParseIP(*localAddress))
	if err != nil {
		log.Fatalf("error constructing server: %v", err)
	}

	probeCh := make(chan struct{})
	resultsCh := make(chan string)

	peerPlugin := &plugin{
		probe:     net.ParseIP(*nlri),
		host:      []byte(*id),
		pingCh:    make(chan ping),
		probeCh:   probeCh,
		resultsCh: resultsCh,
		passive:   *passive,
	}

	err = srv.AddPeer(corebgp.PeerConfig{
		LocalAddress:  net.ParseIP(*localAddress),
		RemoteAddress: net.ParseIP(*remoteAddress),
		LocalAS:       uint32(*localAS),
		RemoteAS:      uint32(*remoteAS),
	}, peerPlugin)

	if err != nil {
		log.Fatalf("error adding peer: %v", err)
	}

	srvErrCh := make(chan error)
	go func() {
		err := srv.Serve([]net.Listener{})
		srvErrCh <- err
	}()

	if *cloudprober {
		go func() {
			serverutils.Serve(func(
				request *epb.ProbeRequest,
				reply *epb.ProbeReply,
			) {
				probeCh <- struct{}{}
				reply.Payload = proto.String(<-resultsCh)
				if err != nil {
					reply.ErrorMessage = proto.String(err.Error())
				}
			})
		}()
	}

	if !*passive {
		go func() {
			update := time.NewTicker(timeInterval)
			for range update.C {
				probeCh <- struct{}{}
			}
		}()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		log.Println("stopping program...")
		srv.Close()
		<-srvErrCh
	case err := <-srvErrCh:
		log.Fatalf("serve error: %v", err)
	}

}

func newMPCap(afi uint16, safi uint8) corebgp.Capability {
	mpData := make([]byte, 4)
	binary.BigEndian.PutUint16(mpData, afi)
	mpData[3] = safi
	return corebgp.Capability{
		Code:  1,
		Value: mpData,
	}
}

func parseType42(attr bgp.PathAttributeInterface) (from []byte, to []byte, ts time.Time, err error) {
	type42 := attr.(*bgp.PathAttributeUnknown)
	if len(type42.Value) < type42Len {
		return nil, nil, time.Now(), fmt.Errorf("incorrect type42 len %d", len(type42.Value))
	}

	decoded := binary.BigEndian.Uint64(type42.Value[30:47])
	ts = time.Unix(int64(decoded), 0)
	return type42.Value[:15], type42.Value[15:30], ts, nil
}

func buildPayload(pingReq ping) []byte {
	result := make([]byte, type42Len)
	copy(result[:15], pingReq.source[:])
	copy(result[15:30], pingReq.dest)
	binary.BigEndian.PutUint64(result[30:], uint64(pingReq.ts))
	return result
}
