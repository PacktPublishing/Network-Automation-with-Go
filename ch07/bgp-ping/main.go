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
	"syscall"
	"time"

	"github.com/jwhited/corebgp"
	bgp "github.com/osrg/gobgp/v3/pkg/packet/bgp"
)

const (
	bgpType42   = 42
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
	passive       = flag.Bool("p", true, "passive mode")
)

type plugin struct {
	localAddr net.IP
	probe     net.IP
	localAS   uint32
	host      []byte
	pingCh    chan ping
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

func (p *plugin) OnEstablished(peer corebgp.PeerConfig, writer corebgp.UpdateMessageWriter) corebgp.UpdateMessageHandler {
	log.Println("peer established")

	log.Printf("Starting main loop")
	go func() {

		period := time.Second * 10
		withdrawAfter := time.Second * 7
		update := time.NewTicker(period)
		withdraw := time.NewTicker(period)
		withdraw.Stop()

		for {
			select {
			case pingReq := <-p.pingCh:
				src := string(bytes.Trim(pingReq.source, "\x00"))

				pingReq.dest = p.host

				type42PathAttr := bgp.NewPathAttributeUnknown(type42Flags, bgpType42, buildPayload(pingReq))

				bytes, err := p.buildUpdate(type42PathAttr)
				if err != nil {
					log.Fatal(err)
				}

				writer.WriteUpdate(bytes)

				log.Printf("sent ping response to %s", src)

			case <-update.C:

				if p.passive {
					continue
				}

				log.Printf("Sending periodic ping")
				pingReq := ping{
					source: p.host,
					ts:     time.Now().Unix(),
				}

				type42PathAttr := bgp.NewPathAttributeUnknown(type42Flags, bgpType42, buildPayload(pingReq))

				bytes, err := p.buildUpdate(type42PathAttr)
				if err != nil {
					log.Fatal(err)
				}

				writer.WriteUpdate(bytes)

				withdraw.Reset(withdrawAfter)

			case <-withdraw.C:
				log.Printf("Sending periodic withdraw")
				bytes, err := p.buildWithdraw()
				if err != nil {
					log.Fatal(err)
				}

				writer.WriteUpdate(bytes)
			}
		}
	}()

	return p.handleUpdate
}

func (p *plugin) buildWithdraw() ([]byte, error) {

	myNLRI := bgp.NewIPAddrPrefix(32, p.probe.String())

	withdrawnRoutes := []*bgp.IPAddrPrefix{myNLRI}

	msg := bgp.NewBGPUpdateMessage(withdrawnRoutes, []bgp.PathAttributeInterface{}, nil)

	return msg.Body.Serialize()
}

func (p *plugin) buildUpdate(type42 *bgp.PathAttributeUnknown) ([]byte, error) {

	withdrawnRoutes := []*bgp.IPAddrPrefix{}

	nexthop := bgp.NewPathAttributeNextHop(p.localAddr.String())

	asPath := bgp.NewPathAttributeAsPath([]bgp.AsPathParamInterface{bgp.NewAs4PathParam(bgp.BGP_ASPATH_ATTR_TYPE_SEQ, []uint32{p.localAS})})

	origin := bgp.NewPathAttributeOrigin(2) // origin incomplete

	myNLRI := bgp.NewIPAddrPrefix(32, p.probe.String())

	msg := bgp.NewBGPUpdateMessage(
		withdrawnRoutes,
		[]bgp.PathAttributeInterface{type42, nexthop, asPath, origin},
		[]*bgp.IPAddrPrefix{myNLRI},
	)

	return msg.Body.Serialize()
}

func (p *plugin) OnClose(peer corebgp.PeerConfig) {
	log.Println("peer closed")
}

func (p *plugin) handleUpdate(peer corebgp.PeerConfig, u []byte) *corebgp.Notification {

	msg, err := bgp.ParseBGPBody(&bgp.BGPHeader{Type: bgp.BGP_MSG_UPDATE}, u)
	if err != nil {
		log.Fatal("failed to parse bgp message ", err)
	}

	if err := bgp.ValidateBGPMessage(msg); err != nil {
		log.Fatal("validate BGP message ", err)
	}

	for _, attr := range msg.Body.(*bgp.BGPUpdate).PathAttributes {

		if attr.GetType() != bgpType42 {
			continue
		}

		source, dest, ts, err := parseType42(attr)
		if err != nil {
			log.Fatal(err)
		}

		if string(bytes.Trim(source, "\x00")) == *id {
			log.Printf("Received a ping response")

			if len(string(bytes.Trim(dest, "\x00"))) == 0 {
				log.Printf("Received a looped response, ignoring")
				continue
			}

			fmt.Printf("bgp_ping_rtt_ms{device=%s} %f", dest, float64(time.Since(ts).Nanoseconds())/1e6)
			return nil
		}

		p.pingCh <- ping{source: source, ts: ts.Unix()}
	}

	return nil
}

func main() {

	flag.Parse()

	log.Print("starting corebgp server")
	srv, err := corebgp.NewServer(net.ParseIP(*localAddress))
	if err != nil {
		log.Fatalf("error constructing server: %v", err)
	}

	p := &plugin{
		localAddr: net.ParseIP(*localAddress),
		probe:     net.ParseIP(*nlri),
		localAS:   uint32(*localAS),
		host:      []byte(*id),
		pingCh:    make(chan ping),
		passive:   *passive,
	}

	err = srv.AddPeer(corebgp.PeerConfig{
		LocalAddress:  net.ParseIP(*localAddress),
		RemoteAddress: net.ParseIP(*remoteAddress),
		LocalAS:       uint32(*localAS),
		RemoteAS:      uint32(*remoteAS),
	}, p)

	if err != nil {
		log.Fatalf("error adding peer: %v", err)
	}

	srvErrCh := make(chan error)
	go func() {
		err := srv.Serve([]net.Listener{})
		srvErrCh <- err
	}()

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
